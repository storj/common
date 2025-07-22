// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information

package combiner_test

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2/combiner"
	"storj.io/common/testcontext"
)

func TestCombiner(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	const n = 1000
	var total int64

	q := combiner.New(ctx, combiner.Options[int64]{
		Process: func(ctx context.Context, queue *combiner.Queue[int64]) {
			for batch := range queue.Batches() {
				for _, v := range batch {
					atomic.AddInt64(&total, v)
				}
			}
		},
		Fail: func(ctx context.Context, queue *combiner.Queue[int64]) {
			t.Fatal("fail should not happen")
		},
		QueueSize: 3,
	})

	var expect int64
	var g errgroup.Group
	for v := range n {
		expect += int64(v)
		g.Go(func() error {
			q.Enqueue(ctx, int64(v))
			return nil
		})
	}
	require.NoError(t, g.Wait())

	require.NoError(t, q.Wait(ctx))
	q.Close()

	require.Equal(t, expect, total)
}

func TestCombiner_Stop(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	const n = 100

	var mu sync.Mutex
	var failed []int64

	q := combiner.New(ctx, combiner.Options[int64]{
		Process: func(ctx context.Context, queue *combiner.Queue[int64]) {
			t.Fatal("process should not happen")
		},
		Fail: func(ctx context.Context, queue *combiner.Queue[int64]) {
			for batch := range queue.Batches() {
				mu.Lock()
				failed = append(failed, batch...)
				mu.Unlock()
			}
		},
		QueueSize: 3,
	})

	// after stop we should not start new jobs
	q.Stop()

	var expect []int64
	var g errgroup.Group
	for v := range n {
		expect = append(expect, int64(v))
		g.Go(func() error {
			q.Enqueue(ctx, int64(v))
			return nil
		})
	}
	require.NoError(t, g.Wait())

	require.NoError(t, q.Wait(ctx))
	q.Close()

	slices.Sort(failed)
	require.Equal(t, expect, failed)
}
