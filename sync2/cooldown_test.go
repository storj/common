// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestCooldown_Basic(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)

	cooldown := sync2.NewCooldown(10 * time.Second)
	defer cooldown.Close()

	count := int64(0)

	completed := make(chan struct{})

	var group errgroup.Group
	cooldown.Start(ctx, &group, func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		completed <- struct{}{}
		return nil
	})

	// make sure cooldown is initialized before running test
	cooldown.Trigger()
	<-completed

	group.Go(func() error {
		defer cooldown.Stop()
		for i := 0; i < 10; i++ {
			cooldown.Trigger()
		}
		return nil
	})

	err := group.Wait()
	require.NoError(t, err)

	endCount := atomic.LoadInt64(&count)
	require.Equal(t, int64(1), endCount)
}

func TestCooldown_MultipleStops(t *testing.T) {
	t.Parallel()

	cooldown := sync2.NewCooldown(time.Second)
	defer cooldown.Close()

	ctx := testcontext.New(t)

	var group errgroup.Group
	var count int64
	cooldown.Start(ctx, &group, func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	})

	go cooldown.Stop()
	cooldown.Stop()
	cooldown.Stop()
}

func TestCooldown_StopCancelled(t *testing.T) {
	t.Parallel()

	cooldown := sync2.NewCooldown(time.Second)
	defer cooldown.Close()

	testCtx := testcontext.New(t)
	ctx, cancel := context.WithCancel(testCtx)
	cancel()

	var group errgroup.Group
	cooldown.Start(ctx, &group, func(_ context.Context) error {
		return nil
	})

	cooldown.Stop()
	cooldown.Stop()
}

func TestCooldown_Stop_EnsureLoopIsFinished(t *testing.T) {
	t.Parallel()

	cooldown := sync2.NewCooldown(time.Second)
	defer cooldown.Close()

	ctx := testcontext.New(t)

	var completed int64
	started := make(chan int)

	go func() {
		_ = cooldown.Run(ctx, func(_ context.Context) error {
			close(started)
			time.Sleep(1 * time.Second)
			atomic.StoreInt64(&completed, 1)
			return nil
		})
	}()

	cooldown.Trigger()
	<-started
	cooldown.Stop()

	require.Equal(t, atomic.LoadInt64(&completed), int64(1))
}
