// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package rpccache

import (
	"container/list"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
)

// TestCache_Expiration checks that inserted entries expire eventually.
func TestCache_Expiration(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		Expiration: time.Nanosecond,
		Close: func(val any) error {
			require.Equal(t, val, "val")
			close(called)
			return nil
		},
	})

	c.Put("key", "val")

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

// TestCache_Expiration_Evicted checks that evicted entries are closed
// even if they have an expiration.
func TestCache_Expiration_Evicted(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		Capacity:   1,
		Expiration: time.Hour,
		Close: func(val any) error {
			require.Equal(t, val, "val0")
			close(called)
			return nil
		},
	})

	c.Put("key0", "val0")
	c.Put("key1", "val1")

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

// TestCache_Stale checks that the stale predicate is called on Take.
func TestCache_Stale(t *testing.T) {
	c := New(Options{
		Stale: func(val any) bool {
			return val == "val0"
		},
	})

	c.Put("key", "val0")
	require.Nil(t, c.Take("key"))
	c.Put("key", "val1")
	require.Equal(t, c.Take("key"), "val1")
}

// TestCache_Capacity checks that total capacity limits are enforced.
func TestCache_Capacity(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		Capacity: 1,
		Close: func(val any) error {
			require.Equal(t, val, "val0")
			close(called)
			return nil
		},
	})

	c.Put("key0", "val0")
	c.Put("key1", "val1") // evicts val0

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

// TestCache_Capacity_Negative checks that negative capacities cache nothing.
func TestCache_Capacity_Negative(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		Capacity: -1,
		Close: func(val any) error {
			require.Equal(t, val, "val")
			close(called)
			return nil
		},
	})

	c.Put("key", "val")

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

// TestCache_KeyCapacity checks that per-key capacity limits are enforced.
func TestCache_KeyCapacity(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		KeyCapacity: 1,
		Close: func(val any) error {
			require.Equal(t, val, "val0")
			close(called)
			return nil
		},
	})

	c.Put("key0", "val0")
	c.Put("key1", "val1")
	c.Put("key0", "val2") // evicts val0

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

// TestCache_KeyCapacity_Negative checks that negative per-key capacities cache nothing.
func TestCache_KeyCapacity_Negative(t *testing.T) {
	ctx := testcontext.New(t)

	called := make(chan struct{})

	c := New(Options{
		KeyCapacity: -1,
		Close: func(val any) error {
			require.Equal(t, val, "val")
			close(called)
			return nil
		},
	})

	c.Put("key", "val")

	select {
	case <-called:
	case <-ctx.Done():
		t.FailNow()
	}
}

func TestCache_ShortExpirationEventuallyClears(t *testing.T) {
	var chanMu sync.Mutex
	var chans []chan struct{}

	newChan := func() chan struct{} {
		ch := make(chan struct{})
		chanMu.Lock()
		chans = append(chans, ch)
		chanMu.Unlock()
		return ch
	}

	defer func() {
		for _, ch := range chans {
			select {
			case <-ch:
			default:
				t.Error("some channel did not close")
				return
			}
		}
	}()

	c := New(Options{
		Expiration: 1,

		Close: func(val any) error {
			go close(val.(chan struct{}))
			return nil
		},

		Stale: func(val any) bool {
			select {
			case <-val.(chan struct{}):
				return true
			default:
				return false
			}
		},
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				x := c.Take("key")
				if x == nil {
					x = newChan()
				}
				c.Put("key", x)
			}
		}()
	}

	wg.Wait()

	start := time.Now()
	for time.Since(start) < 2*time.Second {
		if c.Cached() == 0 {
			return
		}
		runtime.Gosched()
	}

	t.Fatal("cache did not clear")
}

func TestCache_Fuzz(t *testing.T) {
	// set up a unique random state for this test run
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	t.Log("seed:", seed)

	t.Run("Unlimited", func(t *testing.T) { runFuzz(t, rng, 0, 0) })
	t.Run("NoKeyCapacity", func(t *testing.T) { runFuzz(t, rng, 0, -1) })
	t.Run("NoCapacity", func(t *testing.T) { runFuzz(t, rng, -1, 0) })

	sizes := []int{1, 2, 3, 4}
	for _, capSize := range sizes {
		capSize := capSize
		for _, keyCapSize := range sizes {
			keyCapSize := keyCapSize

			title := fmt.Sprintf("Cap:%2d KeyCap:%d", 10*capSize, keyCapSize)
			t.Run(title, func(t *testing.T) { runFuzz(t, rng, capSize, keyCapSize) })
		}
	}

}

func runFuzz(t *testing.T, rng *rand.Rand, capacity, keyCapacity int) {
	// event is some event that happened with the cache
	type event struct {
		key    string
		val    any
		action string // "put" | "take" | "closed"
	}

	// getEvent creates a random event that is used against the cache
	nonce := 0
	getEvent := func() event {
		nonce++
		key := fmt.Sprintf("key%02d", rng.Intn(1e2))
		return event{
			key:    key,
			val:    fmt.Sprintf("%s:val%02d:%d", key, rng.Intn(1e2), nonce),
			action: [2]string{"put", "take"}[rng.Intn(2)],
		}
	}

	// stale defines about 10% of the values to be stale
	stale := func(val any) bool { return val.(string)[6:] >= "val90" }

	// filter removes a value from a slice
	filter := func(vals []any, val any) []any {
		j := 0
		for i := range vals {
			if vals[i] == val {
				continue
			}
			vals[j] = vals[i]
			j++
		}
		return vals[:j]
	}

	// define the log of events
	var log []event

	// create the cache options for the test
	c := New(Options{
		Capacity:    capacity,
		KeyCapacity: keyCapacity,
		Stale:       stale,
		Close: func(val any) error {
			log = append(log, event{
				key:    val.(string)[:5],
				val:    val.(string),
				action: "closed",
			})
			return nil
		},
	})

	// generate the fuzz events
	for i := 0; i < 10000; i++ {
		ev := getEvent()
		switch ev.action {
		case "put":
			c.Put(ev.key, ev.val)
		case "take":
			ev.val = c.Take(ev.key)
		}
		log = append(log, ev)
	}

	// keep track of the events before Close is called
	beforeClose := len(log)
	require.NoError(t, c.Close())

	// check the consistency of the log:
	//   1. any value >= val90 should be closed before a put
	//   2. all of the key capacities are enforced
	//   3. all of the overall capacities are enforced
	//   4. every close call is explained by one of the above
	//   5. every value is eventually closed
	//   6. no values remain in the cache

	state := make(map[string][]any)
	order := make([]any, 0)
	checked := make(map[int]bool)
	openValues := make(map[any]bool)

	// we pre-declare the variables here so that we can get logging of the events
	// but only if the test fails.
	var i int
	var ev event
	defer func() {
		if t.Failed() {
			for i, ev := range log[:i+1] {
				t.Logf("%-5d %+v", i, ev)
			}
		}
	}()

	for i, ev = range log {
		// record that any non-nil value we see is potentially open
		if _, ok := openValues[ev.val]; ev.val != nil && !ok {
			openValues[ev.val] = true
		}

		switch ev.action {
		case "put":
			// property 1: if the value is stale or either capacity is negative
			// then the previous message must be a close for this entry.
			if stale(ev.val) || c.opts.Capacity < 0 || c.opts.KeyCapacity < 0 {
				require.Equal(t, event{
					key:    ev.key,
					val:    ev.val,
					action: "closed",
				}, log[i-1], "event %d", i)
				checked[i-1] = true

				openValues[ev.val] = false
				break
			}

			// add the value to the key and order
			state[ev.key] = append(state[ev.key], ev.val)
			order = append(order, ev.val)

			// property 2: if it puts it over the capacity, the previous
			// message must be a close for this entry.
			if c.opts.KeyCapacity > 0 && len(state[ev.key]) > c.opts.KeyCapacity {
				val := state[ev.key][0]

				require.Equal(t, event{
					key:    ev.key,
					val:    val,
					action: "closed",
				}, log[i-1], "event %d", i)
				checked[i-1] = true

				// remove the entry from the order and per-key state
				openValues[val] = false
				order = filter(order, val)
				state[ev.key] = state[ev.key][1:]
			}

			// property 3: if we're over the overall capacity, the previous
			// message must be a close for the oldest entry.
			if c.opts.Capacity > 0 && len(order) > c.opts.Capacity {
				val := order[0]
				key := val.(string)[:5]

				require.Equal(t, event{
					key:    key,
					val:    val,
					action: "closed",
				}, log[i-1], "event %d", i)
				checked[i-1] = true

				// remove the entry from the order and per-key state
				openValues[val] = false
				order = order[1:]
				state[key] = filter(state[key], val)
			}

		case "take":
			// taking a non-empty key removes it.
			if values := state[ev.key]; len(values) > 0 {
				val := values[len(values)-1]

				openValues[val] = false
				order = filter(order, val)
				state[ev.key] = values[:len(values)-1]
			}

		case "closed":
			// all closes before the cache Close call need to be
			// checked by the other conditions.
			if i < beforeClose {
				continue
			}

			checked[i] = true
			openValues[ev.val] = false
			order = filter(order, ev.val)
			state[ev.key] = filter(state[ev.key], ev.val)
		}
	}

	// check property 4: every closed event is accounted for
	for i, ev := range log {
		require.True(t, ev.action != "closed" || checked[i],
			"event %+v not checked for", ev)
	}

	// check property 5: every value is closed
	for val, open := range openValues {
		require.False(t, open, "value %q was not closed", val)
	}

	// check property 6: no values remain the cache
	require.Len(t, order, 0, "entries left in the age order list")
	for key, values := range state {
		require.Len(t, values, 0, "entries left for key %q", key)
	}
}

//
// benchmarks
//

func runBenchmarkFilterSlice(b *testing.B, n int) {
	ents := make([]*entry, n)
	for i := range ents {
		ents[i] = new(entry)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ent := ents[0]
		ents = filterEntry(ents, ent)
		ents = append(ents, ent) //nolint:makezero // the test removes from slice and adds it back
	}
}

func BenchmarkFilterSlice(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		size := size
		b.Run(strconv.Itoa(size), func(b *testing.B) { runBenchmarkFilterSlice(b, size) })
	}
}

func runBenchmarkFilterList(b *testing.B, n int) {
	l := list.New()
	for i := 0; i < n; i++ {
		l.PushBack(new(entry))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ele := l.Front()
		l.Remove(ele)
		l.PushBack(ele.Value)
	}
}

func BenchmarkFilterList(b *testing.B) {
	for _, size := range []int{10, 100, 1000} {
		size := size
		b.Run(strconv.Itoa(size), func(b *testing.B) { runBenchmarkFilterList(b, size) })
	}
}
