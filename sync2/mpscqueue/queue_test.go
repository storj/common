// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information

package mpscqueue_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2"
	"storj.io/common/sync2/mpscqueue"
)

func TestBasic(t *testing.T) {
	var queue mpscqueue.Queue[int]
	queue.Init()

	const N = 10

	for i := 0; i < N; i++ {
		queue.Enqueue(i)
	}

	for i := 0; i < N; i++ {
		v, ok := queue.Dequeue()
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	_, ok := queue.Dequeue()
	require.False(t, ok)
}

func TestNoBreaks(t *testing.T) {
	var queue mpscqueue.Queue[int]
	queue.Init()

	const N = 10

	var senders errgroup.Group
	for i := 0; i < N; i++ {
		i := i
		senders.Go(func() error {
			queue.Enqueue(i)
			return nil
		})
	}

	// When all senders have completed,
	// then we should be able to read all the values,
	// without failures.
	require.NoError(t, senders.Wait())

	var seen [N]bool
	for {
		v, ok := queue.Dequeue()
		if !ok {
			break
		}
		if seen[v] {
			t.Fatal("encountered value twice", v)
		}
		seen[v] = true
	}

	for i, v := range seen {
		if v == false {
			t.Fatal("did not receive", i)
		}
	}
}

func TestConcurrent(t *testing.T) {
	var queue mpscqueue.Queue[int]
	queue.Init()

	const N = 100

	var seen [N * 3]bool

	_ = sync2.Concurrently(
		func() error {
			for i := 0; i < N; i++ {
				queue.Enqueue(i)
			}
			return nil
		},
		func() error {
			for i := 0; i < N; i++ {
				queue.Enqueue(N + i)
			}
			return nil
		},
		func() error {
			for i := 0; i < N; i++ {
				queue.Enqueue(2*N + i)
			}
			return nil
		},
		func() error {
			seenCount := 0
			for seenCount < 3*N {
				v, ok := queue.Dequeue()
				if !ok {
					runtime.Gosched()
					continue
				}
				if seen[v] {
					return fmt.Errorf("value %v seen already", v)
				}
				seen[v] = true
				seenCount++
			}
			return nil
		})

	for i, v := range seen {
		if !v {
			t.Fatal("did not see", i)
		}
	}
}
