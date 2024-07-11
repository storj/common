// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"context"
	"errors"
	"runtime/trace"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/pb"
	"storj.io/common/sync2"
	"storj.io/drpc"
)

func TestRemoteDebugServer_Unauthenticated(t *testing.T) {
	endpoint := NewEndpoint(func(ctx context.Context) error {
		return errors.New("unauthenticated")
	})

	stream := newFakeStream(context.Background())
	require.Error(t, endpoint.CollectRuntimeTraces(nil, stream), "unauthenticated")
}

func TestRemoteDebugServer_SomeData(t *testing.T) {
	if trace.IsEnabled() {
		t.Skip("tracing already enabled")
	} else if !traceEnabled {
		t.Skip("tracing not enabled")
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

type fakeStream struct {
	drpc.Stream // expected nil but just implements interface

	ctx     context.Context
	written sync2.Fence
	data    []byte
}

func newFakeStream(ctx context.Context) *fakeStream { return &fakeStream{ctx: ctx} }

func (f *fakeStream) Context() context.Context { return f.ctx }
func (f *fakeStream) Send(m *pb.CollectRuntimeTracesResponse) error {
	f.data = append(f.data, m.Data...)
	if len(f.data) > 0 {
		f.written.Release()
	}
	return nil
}
