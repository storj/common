// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"context"
	"runtime"
	"runtime/trace"

	"github.com/zeebo/errs"

	"storj.io/common/pb"
	"storj.io/common/rpc/rpcstatus"
)

// Endpoint implements a remote debug server.
type Endpoint struct {
	pb.DRPCDebugUnimplementedServer

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

	if err := trace.Start(&streamWriter{stream: stream}); err != nil {
		return rpcstatus.Wrap(rpcstatus.FailedPrecondition, errs.New("trace failed to start: %w", err))
	}
	defer trace.Stop()

	<-stream.Context().Done()

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
