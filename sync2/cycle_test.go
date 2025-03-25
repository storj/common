// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestCycle_Basic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var inplace sync2.Cycle
	inplace.SetInterval(time.Second)

	var pointer = sync2.NewCycle(time.Second)

	for _, cycle := range []*sync2.Cycle{pointer, &inplace} {
		cycle := cycle
		t.Run("", func(t *testing.T) {
			t.Parallel()
			defer cycle.Close()

			count := int64(0)

			var group errgroup.Group

			start := time.Now()

			cycle.Start(ctx, &group, func(ctx context.Context) error {
				atomic.AddInt64(&count, 1)
				return nil
			})

			group.Go(func() error {
				defer cycle.Stop()

				const expected = 10
				cycle.Pause()

				startingCount := atomic.LoadInt64(&count)
				for i := 0; i < expected-1; i++ {
					cycle.Trigger()
				}
				cycle.TriggerWait()
				countAfterTrigger := atomic.LoadInt64(&count)

				change := countAfterTrigger - startingCount
				if expected != change {
					return fmt.Errorf("invalid triggers expected %d got %d", expected, change)
				}

				cycle.Restart()
				time.Sleep(3 * time.Second)

				countAfterRestart := atomic.LoadInt64(&count)
				if countAfterRestart == countAfterTrigger {
					return errors.New("cycle has not restarted")
				}

				return nil
			})

			err := group.Wait()
			if err != nil {
				t.Error(err)
			}

			testDuration := time.Since(start)
			if testDuration > 7*time.Second {
				t.Errorf("test took too long %v, expected approximately 3s", testDuration)
			}

			// shouldn't block
			cycle.Trigger()
		})
	}
}

func TestCycle_MultipleStops(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()

	ctx := context.Background()

	var group errgroup.Group
	var count int64
	cycle.Start(ctx, &group, func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	})

	go cycle.Stop()
	cycle.Stop()
	cycle.Stop()

	require.NoError(t, group.Wait())
}

func TestCycle_StopCancelled(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var group errgroup.Group
	cycle.Start(ctx, &group, func(_ context.Context) error {
		return nil
	})

	time.Sleep(2 * time.Second)

	cycle.Stop()
	cycle.Stop()

	err := group.Wait()
	require.True(t, errors.Is(err, context.Canceled))
}

func TestCycle_Disable(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(-1)
	cycle.ChangeInterval(-1)
	require.Panics(t,
		func() {
			cycle.ChangeInterval(5 * time.Minute)
		},
		"changing interval of a disabled cycle should panic",
	)

	executed := false
	cycle = sync2.NewCycle(-1)
	err := cycle.Run(testcontext.New(t), func(ctx context.Context) error {
		executed = true
		return nil
	})
	require.NoError(t, err)
	require.False(t, executed)

	// no op for disabled cycle
	cycle.Pause()
	cycle.Trigger()
	cycle.Restart()
}

func TestCycle_Stop_NotStarted(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	cycle.Stop()
}

func TestCycle_Close_NotStarted(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	cycle.Close()
}

func TestCycle_Stop_EnsureLoopIsFinished(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()

	ctx := context.Background()

	var completed int64
	started := make(chan int)

	go func() {
		_ = cycle.Run(ctx, func(_ context.Context) error {
			close(started)
			time.Sleep(1 * time.Second)
			atomic.StoreInt64(&completed, 1)
			return nil
		})
	}()

	<-started
	cycle.Stop()

	require.Equal(t, atomic.LoadInt64(&completed), int64(1))
}

func TestCycle_TimeTriggered(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()

	ctx := context.Background()

	var ranOnce sync2.Fence

	var group errgroup.Group
	cycle.Start(ctx, &group, func(ctx context.Context) error {
		defer ranOnce.Release()

		if sync2.IsManuallyTriggeredCycle(ctx) {
			return errors.New("shouldn't be manually triggered")
		}
		return nil
	})

	ranOnce.Wait(ctx)
	cycle.Stop()

	require.NoError(t, group.Wait())
}

func TestCycle_ManuallyTriggered(t *testing.T) {
	t.Parallel()

	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()

	ctx := context.Background()

	var group errgroup.Group

	check := false
	cycle.Start(ctx, &group, func(ctx context.Context) error {
		if check {
			if !sync2.IsManuallyTriggeredCycle(ctx) {
				return errors.New("should be manually triggered")
			}
		}
		return nil
	})
	cycle.Pause()
	cycle.TriggerWait()
	check = true
	cycle.TriggerWait()

	cycle.Stop()
	require.NoError(t, group.Wait())
}

func TestCycle_DelayStart(t *testing.T) {
	t.Parallel()
	start := time.Now()
	cycle := sync2.NewCycle(time.Second)
	defer cycle.Close()
	cycle.SetDelayStart()

	ctx := context.Background()

	var ranOnce sync2.Fence

	var group errgroup.Group
	cycle.Start(ctx, &group, func(ctx context.Context) error {
		defer ranOnce.Release()

		testDuration := time.Since(start)
		if testDuration < time.Second {
			return fmt.Errorf("start was not delayed %v, expected >= 1s", testDuration)
		}
		return nil
	})

	ranOnce.Wait(ctx)
	cycle.Stop()

	require.NoError(t, group.Wait())
}
