// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	ctx := testcontext.NewWithTimeout(t, 30*time.Second)

	var group errgroup.Group
	var event sync2.Event
	var done int32

	group.Go(func() error {
		if !event.Wait(ctx) {
			return errors.New("got false from Wait")
		}
		if atomic.LoadInt32(&done) == 0 {
			return errors.New("event not yet signaled")
		}
		return nil
	})

	// Wait a bit for the goroutine to hit the fence;
	// this ensures that the `Wait` doesn't release early.
	time.Sleep(100 * time.Millisecond)

	for range 3 {
		group.Go(func() error {
			atomic.StoreInt32(&done, 1)
			event.Signal()
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestEventSticky(t *testing.T) {
	t.Parallel()

	ctx := testcontext.NewWithTimeout(t, 30*time.Second)
	var event sync2.Event

	// we should be able to signal multiple times without blocking
	event.Signal()
	event.Signal()

	// signaling should be sticky
	if !event.Wait(ctx) {
		t.Error("got false from Wait")
	}
}

func TestEvent_ContextCancel(t *testing.T) {
	t.Parallel()

	tctx := testcontext.NewWithTimeout(t, 30*time.Second)
	ctx, cancel := context.WithCancel(tctx)

	var group errgroup.Group
	var event sync2.Event

	for range 10 {
		group.Go(func() error {
			if event.Wait(ctx) {
				return errors.New("got true from Wait")
			}
			return nil
		})
	}

	// wait a bit for all goroutines to hit the fence
	time.Sleep(100 * time.Millisecond)

	cancel()

	if err := group.Wait(); err != nil {
		t.Fatal(err)
	}
}
