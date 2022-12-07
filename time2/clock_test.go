// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information

package time2_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
	"storj.io/common/time2"
)

const (
	testDuration = time.Minute

	// Tests should run under this time. Normally it should take
	// sub-millisecond but things can happen under CI/CD load that could
	// make this higher. As long as this is less than testDuration, things
	// should be ok.
	failAfter = 2 * time.Second
)

var (
	testDate = time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC)
)

func TestNow(t *testing.T) {
	t.Run("with real clock", func(t *testing.T) {
		now := time2.Now(context.Background())
		require.NotEqual(t, now.String(), testDate.String())
	})
	t.Run("with override", func(t *testing.T) {
		testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			before := time2.Now(ctx)
			timeMachine.Advance(testDuration)
			after := time2.Now(ctx)

			require.Equal(t, testDate.String(), before.String())
			require.Equal(t, testDuration, after.Sub(before))
		})
	})
}

func TestSince(t *testing.T) {
	t.Run("with real clock", func(t *testing.T) {
		before := time.Now()
		time.Sleep(time.Millisecond * 50)
		since := time2.Since(context.Background(), before)
		require.Greater(t, since, time.Duration(0))
	})
	t.Run("with override", func(t *testing.T) {
		testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			before := time2.Now(ctx)
			timeMachine.Advance(testDuration)
			since := time2.Since(ctx, before)
			require.Equal(t, testDuration, since)
		})
	})
}

func TestTicker(t *testing.T) {
	t.Run("with real clock", func(t *testing.T) {
		ticker := time2.NewTicker(context.Background(), time.Millisecond)
		defer ticker.Stop()
		<-ticker.Chan()
		<-ticker.Chan()
		<-ticker.Chan()
	})
	t.Run("with override", func(t *testing.T) {
		elapsed := testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			now := timeMachine.Now()

			var times []time.Time

			ticker := time2.NewTicker(ctx, testDuration)
			defer ticker.Stop()

			select {
			case <-ticker.Chan():
				t.Fatal("ticker should not be ticked")
			default:
			}

			timeMachine.BlockThenAdvance(1, testDuration/2)
			select {
			case <-ticker.Chan():
				t.Fatal("ticker should not be ticked")
			default:
			}

			timeMachine.BlockThenAdvance(1, testDuration/2)
			times = append(times, <-ticker.Chan())
			timeMachine.BlockThenAdvance(1, testDuration)
			times = append(times, <-ticker.Chan())
			timeMachine.BlockThenAdvance(1, testDuration)
			times = append(times, <-ticker.Chan())

			select {
			case <-ticker.Chan():
				t.Fatal("ticker should not be ticked")
			default:
			}

			require.Equal(t, []time.Time{
				now.Add(1 * testDuration),
				now.Add(2 * testDuration),
				now.Add(3 * testDuration),
			}, times)
		})
		require.Less(t, elapsed, failAfter)
	})
}

func TestTimer(t *testing.T) {
	assertStopAndResetReturnValues := func(t *testing.T, timer time2.Timer) {
		assert.False(t, timer.Stop(), "timer should have already expired")
		assert.False(t, timer.Reset(testDuration), "timer should be stopped at the time of reset")
		assert.True(t, timer.Stop(), "timer should still be active because it was just reset")
		assert.False(t, timer.Reset(testDuration), "timer should not be active because it was just stopped")
		assert.True(t, timer.Reset(testDuration), "timer should be active since it was just reset")

	}
	t.Run("with real clock", func(t *testing.T) {
		timer := time2.NewTimer(context.Background(), time.Millisecond)
		defer timer.Stop()
		<-timer.Chan()
		assert.False(t, timer.Reset(time.Millisecond))
		<-timer.Chan()
		assert.False(t, timer.Reset(time.Millisecond))
		<-timer.Chan()

		assertStopAndResetReturnValues(t, timer)
	})

	t.Run("with override", func(t *testing.T) {
		elapsed := testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			now := timeMachine.Now()

			var times []time.Time

			wait := sync2.Go(func() {
				timer := time2.NewTimer(ctx, testDuration)
				defer timer.Stop()
				times = append(times, <-timer.Chan())
				assert.False(t, timer.Reset(testDuration))
				times = append(times, <-timer.Chan())
				assert.False(t, timer.Reset(testDuration))
				times = append(times, <-timer.Chan())

				assertStopAndResetReturnValues(t, timer)
			})

			timeMachine.BlockThenAdvance(1, testDuration)
			timeMachine.BlockThenAdvance(1, testDuration)
			timeMachine.BlockThenAdvance(1, testDuration)

			wait()

			require.Equal(t, []time.Time{
				now.Add(1 * testDuration),
				now.Add(2 * testDuration),
				now.Add(3 * testDuration),
			}, times)
		})
		require.Less(t, elapsed, failAfter)
	})
}

func TestSleep(t *testing.T) {
	t.Run("with real clock", func(t *testing.T) {
		time2.Sleep(context.Background(), time.Millisecond)
	})
	t.Run("with override", func(t *testing.T) {
		elapsed := testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			defer sync2.Go(func() { timeMachine.BlockThenAdvance(1, testDuration) })()
			time2.Sleep(ctx, testDuration)
		})
		require.Less(t, elapsed, failAfter)
	})
	t.Run("with cancellation", func(t *testing.T) {
		elapsed := testWithOverride(t, func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine) {
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			time2.Sleep(ctx, testDuration)
		})
		require.Less(t, elapsed, failAfter)
	})
}

func testWithOverride(tb testing.TB, fn func(ctx context.Context, tb testing.TB, timeMachine *time2.Machine)) time.Duration {
	ctx, timeMachine := time2.WithNewMachine(context.Background(), time2.WithTimeAt(testDate))
	start := time.Now()
	fn(ctx, tb, timeMachine)
	return time.Since(start)
}
