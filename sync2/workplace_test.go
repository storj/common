// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
)

func TestWorkplace(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	firstCtx, firstCancel := context.WithCancel(ctx)
	defer firstCancel()

	var place sync2.Workplace
	completed := 0

	{ // starts ok
		started := place.Start(firstCtx, 256,
			func(jobTag interface{}) bool {
				t.Fatal("shouldn't be called, because there's no job running")
				return false
			},
			func(ctx context.Context) {
				<-ctx.Done()
				completed++
			})
		require.True(t, started)
	}

	{ // finishes ok
		called := false
		started := place.Start(ctx, 1,
			func(jobTag interface{}) bool {
				called = true
				return false
			},
			func(ctx context.Context) { <-ctx.Done() })
		require.True(t, called)
		require.False(t, started)
	}

	{ // overrides, but waits until completes
		started := place.Start(ctx, 32, func(jobTag interface{}) bool {
			return jobTag == 256
		}, func(ctx context.Context) {
			// if the first one and this is started concurrently,
			// then this will report a data-race.
			completed++
		})
		require.True(t, started)
		firstCancel()
	}

	<-place.Done()
	require.Equal(t, 2, completed)
}

func TestWorkplace_Cancel(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	place := sync2.NewWorkPlace()
	place.Cancel()

	ok := place.Start(ctx, 0, nil, func(ctx context.Context) {
		panic("shouldn't be called")
	})
	require.False(t, ok)
}
