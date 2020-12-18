// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestParentLimiterLimiting(t *testing.T) {
	const n, limit = 1000, 10
	parent := sync2.NewParentLimiter(limit)

	child1 := parent.Child()
	child2 := parent.Child()

	ctx := testcontext.New(t)

	counter := int32(0)
	limitChecker := int32(0)
	for i := 0; i < n; i++ {
		child1.Go(ctx, func() {
			if atomic.AddInt32(&limitChecker, 1) > limit {
				t.Fatal("limit exceeded")
			}
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&limitChecker, -1)
			atomic.AddInt32(&counter, 1)
		})

		child2.Go(ctx, func() {
			if atomic.AddInt32(&limitChecker, 1) > limit {
				t.Fatal("limit exceeded")
			}
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&limitChecker, -1)
			atomic.AddInt32(&counter, 1)
		})
	}

	parent.Wait()
	require.Zero(t, atomic.LoadInt32(&limitChecker))
	require.Equal(t, int32(n*2), atomic.LoadInt32(&counter))
}

func TestChildLimiterCancelling(t *testing.T) {
	const n, limit = 1000, 10
	parent := sync2.NewParentLimiter(limit)

	child1 := parent.Child()
	child2 := parent.Child()

	ctxChild1, cancel := context.WithCancel(testcontext.New(t))
	ctxChild2 := testcontext.New(t)
	defer ctxChild2.Cleanup()

	counterChild1 := int32(0)
	counterChild2 := int32(0)
	block := make(chan struct{})
	allReturned := make(chan struct{})
	go func() {
		cancel()
		for i := 0; i < n; i++ {
			child1.Go(ctxChild1, func() {
				atomic.AddInt32(&counterChild1, 1)
				<-block
			})

			child2.Go(ctxChild2, func() {
				atomic.AddInt32(&counterChild2, 1)
				<-block
			})
		}

		close(allReturned)
	}()

	close(block)
	<-allReturned

	parent.Wait()

	assert.Condition(t, func() bool {
		execTasks := atomic.LoadInt32(&counterChild1)
		return execTasks > 0 && execTasks < n
	}, "child 1 shold not have run all the tasks")

	assert.Equal(t, int32(n), atomic.LoadInt32(&counterChild2))
}

func TestChildLimiter_Wait(t *testing.T) {
	parent := sync2.NewParentLimiter(10)

	child1 := parent.Child()
	child2 := parent.Child()

	ctx := testcontext.New(t)

	waitForFinishing := make(chan struct{})

	child1.Go(ctx, func() {
		waitForFinishing <- struct{}{}
	})

	child2.Go(ctx, func() {
		waitForFinishing <- struct{}{}
	})

	waitsDone := make(chan struct{})
	go func() {
		child1.Wait()
		child2.Wait()
		waitsDone <- struct{}{}
	}()

	counter := int32(0)
	for i := 0; i < 3; i++ {
		select {
		case <-waitsDone:
			require.Equal(t, int32(2), atomic.LoadInt32(&counter))
		case <-waitForFinishing:
			atomic.AddInt32(&counter, 1)
		}
	}
}
