// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestBarrier(t *testing.T) {
	t.Parallel()

	ctx := testcontext.NewWithTimeout(t, 30*time.Second)

	const n = 10
	barrier := sync2.NewBarrier(n)

	var arrived atomic.Int32
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			arrived.Add(1)
			require.NoError(t, barrier.Wait(ctx))
			// nobody passes before everyone has arrived.
			require.EqualValues(t, n, arrived.Load())
		})
	}
	wg.Wait()

	// callers after the release pass through.
	require.NoError(t, barrier.Wait(ctx))
}

func TestBarrier_Canceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(testcontext.New(t))
	cancel()

	barrier := sync2.NewBarrier(2)
	require.ErrorIs(t, barrier.Wait(ctx), context.Canceled)
}

func TestBarrier_Zero(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	require.NoError(t, sync2.NewBarrier(0).Wait(ctx))
}
