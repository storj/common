// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package eventstat

import (
	"sync"
)

type counter struct {
	mu        sync.Mutex
	counters  map[string]uint64
	name      string
	key       string
	discarded bool

	// an empty (but already allocated) map from previous attempt to re-use
	free map[string]uint64

	// number of counters stored in memory
	limit int
}

func (c *counter) publishAndReset(publish Publisher) {
	discardedMarker := float64(0)
	c.mu.Lock()
	counters := c.counters
	if c.free == nil {
		c.counters = map[string]uint64{}
	} else {
		// re-using one of the previous (but empty) map from the memory
		c.counters = c.free
		c.free = nil
	}
	if c.discarded {
		discardedMarker = 1
	}
	c.discarded = false
	c.mu.Unlock()

	for name, count := range counters {
		publish(c.name+"_count", Tags{c.key: name}, float64(count))
	}
	publish(c.name+"_buckets", Tags{}, float64(len(counters)))

	publish(c.name+"_discarded", Tags{}, discardedMarker)

	// clean up the original map, but keep a reference to the memory space to avoid a new allocation
	for k := range counters {
		delete(counters, k)
	}
	c.mu.Lock()
	// this is an empty map, can be used to initialize a new map with using existing memory space
	c.free = counters
	c.mu.Unlock()

}

// Increment bumps the usage count of one of the counters.
func (c *counter) increment(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// safety valve, hard limit the memory / network usage
	if len(c.counters) < c.limit {
		c.counters[name]++
		return
	}

	// no new counters, but bump the value
	_, found := c.counters[name]
	if !found {
		c.counters["<DISCARDED>"]++
		c.discarded = true
		return
	}

	c.counters[name]++
}
