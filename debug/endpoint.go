// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"context"
	"runtime"
	"runtime/trace"
	"sync"
	"sync/atomic"

	"github.com/zeebo/errs"

	"storj.io/common/pb"
	"storj.io/common/rpc/rpcstatus"
)

// Endpoint implements a remote debug server.
type Endpoint struct {
	pb.DRPCDebugUnimplementedServer

	mu sync.Mutex

	Auth func(ctx context.Context) error
}

// NewEndpoint constructs a RemoteDebugEndpoint that will consult the given auth function
// with the request context to determine if the request is authorized.
func NewEndpoint(auth func(ctx context.Context) error) *Endpoint {
	return &Endpoint{Auth: auth}
}

// CollectRuntimeTraces will stream trace data to the client until the client cancels the request
// either explicitly or some error happens in sending.
func (f *Endpoint) CollectRuntimeTraces(_ *pb.CollectRuntimeTracesRequest, stream pb.DRPCDebug_CollectRuntimeTracesStream) error {
	if err := f.Auth(stream.Context()); err != nil {
		return rpcstatus.Wrap(rpcstatus.Unauthenticated, err)
	}
	if !traceEnabled {
		return rpcstatus.Wrap(rpcstatus.FailedPrecondition, errs.New("trace is not enabled: %v", runtime.Version()))
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := trace.Start(&streamWriter{stream: stream}); err != nil {
		return rpcstatus.Wrap(rpcstatus.FailedPrecondition, errs.New("trace failed to start: %w", err))
	}
	defer trace.Stop()

	<-stream.Context().Done()

	return nil
}

// CollectRuntimeTraces2 will stream trace data to the client until the client sends a done message
// some error happens, and it then flushes the trace data and captured packet data.
func (f *Endpoint) CollectRuntimeTraces2(stream pb.DRPCDebug_CollectRuntimeTraces2Stream) error {
	if err := f.Auth(stream.Context()); err != nil {
		return rpcstatus.Wrap(rpcstatus.Unauthenticated, err)
	}
	if !traceEnabled {
		return rpcstatus.Wrap(rpcstatus.FailedPrecondition, errs.New("trace is not enabled: %v", runtime.Version()))
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := trace.Start(&streamWriter{stream: stream}); err != nil {
		return rpcstatus.Wrap(rpcstatus.FailedPrecondition, errs.New("trace failed to start: %w", err))
	}
	stop := new(atomic.Bool)
	defer func() {
		if !stop.Swap(true) {
			trace.Stop()
		}
	}()

	// launch a goroutine to log packets into the trace log and wait for it to exit.
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		capturePackets(stream.Context(), stop)
		wg.Done()
	}()

	// wait for a done message or error from caller
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		if msg.Done {
			break
		}
	}

	stop.Store(true)
	trace.Stop()

	return nil
}

type streamWriter struct {
	stream pb.DRPCDebug_CollectRuntimeTracesStream
}

func (s *streamWriter) Write(p []byte) (int, error) {
	if err := s.stream.Send(&pb.CollectRuntimeTracesResponse{Data: p}); err != nil {
		_ = s.stream.Close()
		return 0, err
	}
	return len(p), nil
}
