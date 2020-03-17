// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
)

func TestNewSuccessThreshold(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		desc             string
		tasks            int
		successThreshold float64
		isError          bool
	}{
		{
			desc:             "OK",
			tasks:            10,
			successThreshold: 0.75,
			isError:          false,
		},
		{
			desc:             "OK",
			tasks:            134,
			successThreshold: 1,
			isError:          false,
		},
		{
			desc:             "Error: invalid tasks (0)",
			tasks:            0,
			successThreshold: 0.75,
			isError:          true,
		},
		{
			desc:             "Error: invalid tasks (1)",
			tasks:            1,
			successThreshold: 0.75,
			isError:          true,
		},
		{
			desc:             "Error: invalid tasks (negative)",
			tasks:            -23,
			successThreshold: 0.75,
			isError:          true,
		},
		{
			desc:             "Error: invalid successThreshold (0)",
			tasks:            134,
			successThreshold: 0,
			isError:          true,
		},
		{
			desc:             "Error: invalid successThreshold (negative)",
			tasks:            134,
			successThreshold: -1.5,
			isError:          true,
		},
		{
			desc:             "Error: invalid successThreshold (greater than 1)",
			tasks:            134,
			successThreshold: 1.00001,
			isError:          true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			_, err := sync2.NewSuccessThreshold(tc.tasks, tc.successThreshold)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSuccessThreshold_AllSuccess(t *testing.T) {
	t.Parallel()

	const (
		tasks           = 10
		threshold       = 0.7
		successfulTasks = 7
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	for i := 0; i < tasks; i++ {
		go func() {
			successThreshold.Success()
			wg.Done()
		}()
	}

	successThreshold.Wait(context.Background())
	wg.Wait()

	require.Equal(t, tasks, successThreshold.SuccessCount())
	require.Equal(t, 0, successThreshold.FailureCount())
}

func TestSuccessThreshold_AllFailures(t *testing.T) {
	t.Parallel()

	const (
		tasks     = 10
		threshold = 0.7
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	for i := 0; i < tasks; i++ {
		go func() {
			successThreshold.Failure()
			wg.Done()
		}()
	}

	successThreshold.Wait(context.Background())
	wg.Wait()

	require.Equal(t, 0, successThreshold.SuccessCount())
	require.Equal(t, tasks, successThreshold.FailureCount())
}

func TestSuccessThreshold_FailuresWithReachedSuccessThreshold(t *testing.T) {
	t.Parallel()

	const (
		tasks           = 10
		threshold       = 0.4
		successfulTasks = 4
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	successfulTasksDone := make(chan struct{}, successfulTasks)
	for i := 0; i < tasks; i++ {
		go func(i int) {

			// Alternate tasks with success & failure
			if i%2 == 0 {
				successfulTasksDone <- struct{}{}
				successThreshold.Success()
			} else {
				successThreshold.Failure()
			}

			wg.Done()
		}(i)
	}

	successThreshold.Wait(context.Background())
	// Check that Wait unblocked when reached the successThreshold
	require.Len(t, successfulTasksDone, cap(successfulTasksDone))

	require.Equal(t, successfulTasks, successThreshold.SuccessCount())

	// purge the rest of the goroutines
	for i := successfulTasks; i < tasks/2; i++ {
		<-successfulTasksDone
	}

	wg.Wait()
}

func TestSuccessThreshold_FailuresWithoutReachedSuccessThreshold(t *testing.T) {
	t.Parallel()

	const (
		tasks     = 10
		threshold = 0.8
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	for i := 0; i < tasks; i++ {
		go func(i int) {
			// Alternate tasks with success & failure
			if i%2 == 0 {
				successThreshold.Success()
			} else {
				successThreshold.Failure()
			}

			wg.Done()
		}(i)
	}

	successThreshold.Wait(context.Background())
	wg.Wait()

	require.Equal(t, tasks/2, successThreshold.SuccessCount())
	require.Equal(t, tasks/2, successThreshold.FailureCount())
}

func TestSuccessThreshold_ExtraTasksAreFine(t *testing.T) {
	t.Parallel()

	const (
		tasks      = 10
		threshold  = 0.7
		extraTasks = 5
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks + extraTasks)
	for i := 0; i < (tasks + extraTasks); i++ {
		go func(i int) {
			if i%2 == 0 {
				successThreshold.Success()
			} else {
				successThreshold.Failure()
			}

			wg.Done()
		}(i)
	}

	successThreshold.Wait(context.Background())
	wg.Wait()
}

func TestSuccessThreshold_SuccessRateCloseTo0(t *testing.T) {
	t.Parallel()

	const (
		tasks             = 2
		threshold         = 0.1
		expectedThreshold = 1
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	completedTasks := make(chan struct{}, expectedThreshold)
	for i := 0; i < tasks; i++ {
		go func() {
			completedTasks <- struct{}{}
			successThreshold.Success()

			wg.Done()
		}()
	}

	successThreshold.Wait(context.Background())
	// Check that Wait unblocked when reached the successThreshold
	require.Len(t, completedTasks, cap(completedTasks))

	// purge the rest of the goroutines
	for i := expectedThreshold; i < tasks; i++ {
		<-completedTasks
	}

	wg.Wait()
}

func TestSuccessThreshold_SuccessThresholdNumTasks(t *testing.T) {
	t.Parallel()

	const (
		tasks             = 2
		threshold         = 1
		expectedThreshold = 2
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	for i := 0; i < tasks; i++ {
		go func() {
			successThreshold.Success()
			wg.Done()
		}()
	}

	successThreshold.Wait(context.Background())
	wg.Wait()
}

func TestSuccessThreshold_CallingWaitMoreThanOnce(t *testing.T) {
	t.Parallel()

	const (
		tasks           = 10
		threshold       = 0.7
		successfulTasks = 7
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(tasks)
	for i := 0; i < tasks; i++ {
		go func() {
			successThreshold.Success()
			wg.Done()
		}()
	}

	successThreshold.Wait(context.Background())

	// These two wait calls must not block
	successThreshold.Wait(context.Background())
	successThreshold.Wait(context.Background())
	wg.Wait()
}

func TestSuccessThreshold_CancellingWait(t *testing.T) {
	t.Parallel()

	const (
		tasks           = 10
		threshold       = 0.7
		successfulTasks = 7
	)

	successThreshold, err := sync2.NewSuccessThreshold(tasks, threshold)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	successThreshold.Wait(ctx)
}
