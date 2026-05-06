// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2_test

import (
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

// TestWithTimeout_Panic confirms that panics that occur within the goroutine
// executing the task callback are surfaced to the caller.
func TestWithTimeout_Panic(t *testing.T) {
	panicValue := "panicked"

	t.Run("Panic before timeout", func(t *testing.T) {
		require.PanicsWithValue(t, panicValue, func() {
			sync2.WithTimeout(time.Second, func() {
				time.Sleep(time.Millisecond * 30)
				panic(panicValue)
			}, func() {})
		})
	})

	t.Run("Panic after timeout", func(t *testing.T) {
		ch := make(chan struct{})
		require.PanicsWithValue(t, panicValue, func() {
			sync2.WithTimeout(time.Second, func() {
				<-ch
				panic(panicValue)
			}, func() {
				close(ch)
			})
		})
	})

	t.Run("Task panic takes precedence over timeout panic", func(t *testing.T) {
		ch := make(chan struct{})
		require.PanicsWithValue(t, panicValue, func() {
			sync2.WithTimeout(time.Second, func() {
				<-ch
				panic(panicValue)
			}, func() {
				close(ch)
				panic("panicked in timeout callback")
			})
		})
	})
}

func BenchmarkWithTimeout(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sync2.WithTimeout(time.Second, func() {}, func() {})
	}
}
