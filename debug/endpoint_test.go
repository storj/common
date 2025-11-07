// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"context"
	"errors"
	"io"
	"runtime/trace"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/pb"
	"storj.io/common/sync2"
	"storj.io/drpc"
)

func TestRemoteDebugServer_CollectRuntimeTraces_Unauthenticated(t *testing.T) {
	endpoint := NewEndpoint(func(ctx context.Context) error {
		return errors.New("unauthenticated")
	})

	stream := newFakeStream(context.Background())
	require.Error(t, endpoint.CollectRuntimeTraces(nil, stream), "unauthenticated")
}

func TestRemoteDebugServer_CollectRuntimeTraces_SomeData(t *testing.T) {
	if trace.IsEnabled() {
		t.Skip("tracing already enabled")
	}

	endpoint := NewEndpoint(func(ctx context.Context) error { return nil })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := newFakeStream(ctx)

	var group errgroup.Group
	group.Go(func() error {
		return endpoint.CollectRuntimeTraces(nil, stream)
	})
	group.Go(func() error {
		<-stream.written.Done()
		cancel()
		return nil
	})
	require.NoError(t, group.Wait())

	require.NotZero(t, len(stream.data))
}

func TestRemoteDebugServer_CollectRuntimeTraces2_Unauthenticated(t *testing.T) {
	endpoint := NewEndpoint(func(ctx context.Context) error {
		return errors.New("unauthenticated")
	})

	stream := newFakeStream(context.Background())
	require.Error(t, endpoint.CollectRuntimeTraces2(stream), "unauthenticated")
}

func TestRemoteDebugServer_CollectRuntimeTraces2_SomeData(t *testing.T) {
	if trace.IsEnabled() {
		t.Skip("tracing already enabled")
	}

	endpoint := NewEndpoint(func(ctx context.Context) error { return nil })

	stream := newFakeStream(context.Background())

	var group errgroup.Group
	group.Go(func() error {
		return endpoint.CollectRuntimeTraces2(stream)
	})
	group.Go(func() error {
		<-stream.written.Done()
		stream.InjectDone(true)
		return nil
	})
	require.NoError(t, group.Wait())

	require.NotZero(t, len(stream.data))
}

type fakeStream struct {
	drpc.Stream // expected nil but just implements interface

	ctx     context.Context
	written sync2.Fence
	data    []byte
	reqs    chan bool
}

func newFakeStream(ctx context.Context) *fakeStream {
	return &fakeStream{
		ctx:  ctx,
		reqs: make(chan bool),
	}
}

func (f *fakeStream) InjectDone(done bool) { f.reqs <- done }

func (f *fakeStream) Context() context.Context { return f.ctx }
func (f *fakeStream) Send(m *pb.CollectRuntimeTracesResponse) error {
	f.data = append(f.data, m.Data...)
	if len(f.data) > 0 {
		f.written.Release()
	}
	return nil
}
func (f *fakeStream) Recv() (*pb.CollectRuntimeTracesRequest, error) {
	done, ok := <-f.reqs
	if !ok {
		return nil, io.EOF
	}
	return &pb.CollectRuntimeTracesRequest{Done: done}, nil
}
