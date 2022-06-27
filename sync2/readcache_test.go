// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/errs2"
	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestReadCache(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	cacheCtx, cacheCancel := context.WithCancel(ctx)
	defer cacheCancel()

	start := time.Now()

	refresh := int64(0)
	cache, err := sync2.NewReadCache(2*time.Second, 5*time.Second,
		func(ctx context.Context) (interface{}, error) {
			return atomic.AddInt64(&refresh, 1), nil
		})
	require.NoError(t, err)

	ctx.Go(func() error { return cache.Run(cacheCtx) })

	{ // first call should trigger a refresh
		result, err := cache.Get(ctx, start)
		require.NoError(t, err)
		require.Equal(t, int64(1), result)
	}

	{ // call at the same time shouldn't trigger a refresh
		result, err := cache.Get(ctx, start)
		require.NoError(t, err)
		require.Equal(t, int64(1), result)
	}

	{ // call at a later time shouldn't trigger a refresh until we reach refresh threshold
		result, err := cache.Get(ctx, start.Add(time.Second))
		require.NoError(t, err)
		require.Equal(t, int64(1), result)
	}

	{ // after 2s we should trigger a background refresh, however, should return the latest result
		result, err := cache.Get(ctx, start.Add(2*time.Second))
		require.NoError(t, err)
		require.Equal(t, int64(1), result)

		updated, err := cache.Wait(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(2), updated)
	}

	{ // after 5s we should trigger a refresh and wait for the result.
		result, err := cache.Get(ctx, start.Add(2*time.Second+5*time.Second))
		require.NoError(t, err)
		require.Equal(t, int64(3), result)
	}
}

func TestReadCache_Cancellation(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	cacheCtx, cacheCancel := context.WithCancel(ctx)
	defer cacheCancel()

	cache, err := sync2.NewReadCache(2*time.Second, 5*time.Second,
		func(ctx context.Context) (interface{}, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})
	require.NoError(t, err)

	ctx.Go(func() error { return cache.Run(cacheCtx) })

	cacheCancel()
	ctx.Wait()

	state, err := cache.Get(ctx, time.Now())
	require.Error(t, err)
	require.True(t, errs2.IsCanceled(err))
	require.Nil(t, state)
}

func TestReadCache_Error(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	cacheCtx, cacheCancel := context.WithCancel(ctx)
	defer cacheCancel()

	result := make(chan int, 3)
	result <- -1

	cache, err := sync2.NewReadCache(2*time.Second, 5*time.Second,
		func(ctx context.Context) (interface{}, error) {
			select {
			case v := <-result:
				if v == -1 {
					return nil, errs.New("failure")
				}
				return v, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		})
	require.NoError(t, err)

	ctx.Go(func() error { return cache.Run(cacheCtx) })

	start := time.Now()

	// read an error
	state, err := cache.Get(ctx, start)
	require.Error(t, err)
	require.Nil(t, state)

	// when cache gets the result early, there is going to be a write race.
	testRace1, testRace2 := 0, 0

	var wg sync.WaitGroup
	wg.Add(2)

	// concurrent read, should trigger refresh
	ctx.Go(func() error {
		defer wg.Done()
		value, err := cache.Get(ctx, start)
		testRace1 = 1
		if value != 1 {
			return errs.New("wrong result")
		}
		return err
	})

	// concurrent read, not triggering refresh
	ctx.Go(func() error {
		defer wg.Done()

		value, err := cache.Get(ctx, start)
		testRace2 = 1
		if value != 1 {
			return errs.New("wrong result")
		}
		return err
	})

	testRace1 = 0
	testRace2 = 0

	// both reads should get this result
	result <- 1
	// this should be ignored
	result <- -1

	wg.Wait()

	t.Log(testRace1, testRace2)
}

func TestReadCache_Concurrent(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	cacheCtx, cacheCancel := context.WithCancel(ctx)
	defer cacheCancel()

	result := make(chan int, 1)
	result <- 1

	cache, err := sync2.NewReadCache(2*time.Second, 5*time.Second,
		func(ctx context.Context) (interface{}, error) {
			select {
			case v := <-result:
				return v, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		})
	require.NoError(t, err)

	ctx.Go(func() error { return cache.Run(cacheCtx) })

	// cache some result
	start := time.Now()
	value, err := cache.Get(ctx, start)
	require.NoError(t, err)
	require.Equal(t, 1, value)

	var wg sync.WaitGroup
	wg.Add(2)

	// concurrent read, not triggering refresh
	ctx.Go(func() error {
		defer wg.Done()
		value, err := cache.Get(ctx, start)
		if value != 1 {
			return errs.New("wrong result")
		}
		return err
	})

	// concurrent read, triggering a refresh, but not waiting
	ctx.Go(func() error {
		defer wg.Done()
		value, err := cache.Get(ctx, start.Add(time.Second*3))
		if value != 1 {
			return errs.New("wrong result")
		}
		return err
	})

	// concurrent read, triggering a refresh, but waiting
	ctx.Go(func() error {
		_, err := cache.Get(ctx, start.Add(time.Second*9))
		if err == nil {
			return errs.New("did not get an error due to cancellation")
		}
		return nil
	})

	wg.Wait()
	cacheCancel()
}
