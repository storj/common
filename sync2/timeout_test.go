// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestWithTimeout(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	workResult := 0
	timeoutResult := 0

	ctx.Go(func() error {
		sync2.WithTimeout(time.Second, func() {
			workResult = 1
		}, func() {
			timeoutResult = 1
		})
		return nil
	})

	ctx.Wait()

	require.Equal(t, workResult, 1)
	require.Equal(t, timeoutResult, 0)
}

func TestWithTimeout_Fail(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	working := make(chan struct{})

	workResult := 0
	timeoutResult := 0
	allResult := 0

	workCompleted := make(chan struct{})
	timeoutCompleted := make(chan struct{})

	ctx.Go(func() error {
		sync2.WithTimeout(0, func() {
			<-working
			workResult = 1
			close(workCompleted)
		}, func() {
			timeoutResult = 1
			close(timeoutCompleted)
		})
		allResult = 1
		return nil
	})

	<-timeoutCompleted

	require.Equal(t, workResult, 0)
	require.Equal(t, timeoutResult, 1)
	require.Equal(t, allResult, 0)

	close(working)
	<-workCompleted

	require.Equal(t, workResult, 1)
	require.Equal(t, timeoutResult, 1)
}

func TestWithTimeout_NoOnTimeoutAfterDo(t *testing.T) {
	t.Parallel()

	// Force the race window by making do's completion collide with the timer's
	// deadline: do sleeps for `timeout` so it returns at roughly the same
	// instant the timer fires. With the pre-CAS implementation, if the timer
	// expires between do returning and the deferred t.Stop(), onTimeout runs
	// spuriously after do has already completed (this is the piecestore bug:
	// a successful Send followed by the timeout callback cancelling the
	// shared ctx). The fix must ensure onTimeout never observes a completed
	// do.
	const timeout = time.Millisecond
	const iterations = 200
	var spurious atomic.Int32
	for i := 0; i < iterations; i++ {
		var doDone atomic.Bool
		sync2.WithTimeout(timeout, func() {
			time.Sleep(timeout)
			doDone.Store(true)
		}, func() {
			if doDone.Load() {
				spurious.Add(1)
			}
		})
	}
	require.Zero(t, spurious.Load(),
		"onTimeout ran after do completed in %d/%d iterations",
		spurious.Load(), iterations)
}

func BenchmarkWithTimeout(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sync2.WithTimeout(time.Second, func() {}, func() {})
	}
}
