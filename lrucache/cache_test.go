// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package lrucache

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/testcontext"
)

func TestCache_LRU(t *testing.T) {
	cache := NewOf[string](Options{Capacity: 2})
	check := newChecker(t, cache)

	check("a", 1)
	check("a", 1)
	check("b", 2)
	check("a", 2)
	check("c", 3)
	check("b", 4)
	check("c", 4)
	check("a", 5)
}

func TestCache_Expires(t *testing.T) {
	cache := NewOf[string](Options{Capacity: 2, Expiration: time.Nanosecond})
	check := newChecker(t, cache)

	check("a", 1)
	time.Sleep(time.Second)
	check("a", 2)
}

func TestCache_Get_Fuzz(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	cache := New(Options{Capacity: 2, Expiration: 100 * time.Millisecond})
	keys := "abcdefghij"

	var ops uint64
	procs := runtime.GOMAXPROCS(-1)

	for i := 0; i < procs; i++ {
		ctx.Go(func() error {
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			for {
				if atomic.AddUint64(&ops, 1) > 1000000 {
					return nil
				}

				shouldErr := rng.Intn(10) == 0
				ran := false
				kidx := rng.Intn(len(keys))
				key := keys[kidx : kidx+1]

				value, err := cache.Get(key, func() (interface{}, error) {
					ran = true
					if shouldErr {
						return nil, errs.New("random error")
					}
					return key, nil
				})

				if ran {
					if shouldErr && err == nil {
						return errs.New("should have errored and did not")
					}
					if !shouldErr && err != nil {
						return errs.New("should not have errored but did")
					}
				}
				if value != key && !(ran && shouldErr) {
					return errs.New("expected %q but got %q", key, value)
				}
			}
		})
	}

	ctx.Wait()
}

func TestCache_Get_Dedup(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	cache := New(Options{Capacity: 1})
	fnCalled := make(chan struct{})

	ctx.Go(func() error {
		_, _ = cache.Get("key", func() (interface{}, error) {
			fnCalled <- struct{}{}
			time.Sleep(time.Millisecond * 10)
			return 1, nil
		})

		return nil
	})

	<-fnCalled

	value, err := cache.Get("key", func() (interface{}, error) {
		return 0, nil
	})

	require.NoError(t, err)
	require.Equal(t, 1, value)
}

func TestCache_Add_and_GetCached(t *testing.T) {
	cache := New(Options{Capacity: 2, Expiration: time.Millisecond})

	// Never added.
	_, cached := cache.GetCached("key1")
	require.False(t, cached, "GetCached -> cached")

	// Never added before.
	replaced := cache.Add("key1", 1)
	require.False(t, replaced, "Add -> replaced")
	value, cached := cache.GetCached("key1")
	require.True(t, cached, "GetCached -> cached")
	require.Equal(t, 1, value)
	// Added before.
	replaced = cache.Add("key1", 1)
	require.True(t, replaced, "Add -> replaced")

	// Added before but expired.
	time.Sleep(time.Millisecond)
	replaced = cache.Add("key1", 1)
	require.False(t, replaced, "Add -> replaced")

	// Never added before.
	replaced = cache.Add("key2", 2)
	require.False(t, replaced, "Add -> replaced")
	replaced = cache.Add("key3", 3)
	require.False(t, replaced, "Add -> replaced")

	// Evicted because of capacity limit.
	_, cached = cache.GetCached("key1")
	require.False(t, cached, "GetCached -> cached (evicted because it was the least recently used)")
	value, cached = cache.GetCached("key2")
	require.True(t, cached, "GetCached -> cached")
	require.Equal(t, 2, value)
	value, cached = cache.GetCached("key3")
	require.True(t, cached, "GetCached -> cached")
	require.Equal(t, 3, value)
}

func TestCache_Add_and_GetCached_Fuzz(t *testing.T) {
	const numEntries = 200
	require.Zero(t, numEntries%2) // Ensure that numEntries is even.

	ctx := testcontext.New(t)
	cache := NewOf[int64](Options{Capacity: numEntries})

	// Spin up 2 Goroutines that add values and counts the added elements for
	// each one.

	var addCounter1 int64 = -1
	ctx.Go(func() error {
		for e := int64(0); e < numEntries/2; e++ {
			replaced := cache.Add(fmt.Sprintf("%d", e), e)
			atomic.AddInt64(&addCounter1, 1)

			require.False(t, replaced, "replaced")
		}

		return nil
	})

	var addCounter2 int64 = (numEntries / 2) - 1
	ctx.Go(func() error {
		for e := int64(numEntries / 2); e < numEntries; e++ {
			replaced := cache.Add(fmt.Sprintf("%d", e), e)
			atomic.AddInt64(&addCounter2, 1)

			require.False(t, replaced, "replaced")
		}

		return nil
	})

	// Spin up 2 Goroutines for getting values, one uses keys that exist and
	// another one that use keys that don't exist.

	ctx.Go(func() error {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		for e := 0; e < numEntries; {
			expVal := rng.Int63n(numEntries)
			addCounter1 := atomic.LoadInt64(&addCounter1)
			addCounter2 := atomic.LoadInt64(&addCounter2)
			if expVal > addCounter2 || (expVal > addCounter1 && expVal < numEntries/2) {
				// The value isn't in the cache yet.
				continue
			}

			e++
			value, cached := cache.GetCached(fmt.Sprintf("%d", expVal))
			require.True(t, cached, "cached")
			require.Equal(t, expVal, value, "value")
		}

		return nil
	})

	ctx.Go(func() error {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		for e := uint64(0); e < numEntries; e++ {
			key := rng.Int63n(numEntries) + numEntries

			_, cached := cache.GetCached(fmt.Sprintf("%d", key))
			require.False(t, cached, "cached")
		}

		return nil
	})

	ctx.Wait()
}

//
// helper
//

type checker struct {
	t     *testing.T
	cache *ExpiringLRUOf[string]
	calls int
}

func newChecker(t *testing.T, cache *ExpiringLRUOf[string]) func(string, int) {
	return (&checker{t: t, cache: cache}).Check
}

func (c *checker) makeCallback(v string) func() (string, error) {
	return func() (string, error) {
		c.calls++
		return v, nil
	}
}

func (c *checker) Check(key string, calls int) {
	value, err := c.cache.Get(key, c.makeCallback(key))
	require.Equal(c.t, c.calls, calls)
	require.Equal(c.t, value, key)
	require.NoError(c.t, err)
}
