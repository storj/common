// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package eventstat

import (
	"sync"
)

// Tags represent key/values for any event.
type Tags map[string]string

// Publisher is a function which sends out statistics.
type Publisher func(name string, tags Tags, value float64)

// Sink is a function to receive an event.
type Sink func(name string)

// Registry represents the collection of different event counters.
type Registry struct {
	mu       sync.RWMutex
	counters []*counter
}

// NewTagCounter creates an event counter which is registered to the registry.
func (r *Registry) NewTagCounter(name string, key string, opts ...func(*counter)) Sink {
	r.mu.Lock()
	defer r.mu.Unlock()
	e := counter{
		name:     name,
		counters: make(map[string]uint64),
		key:      key,
		limit:    1000,
	}
	for _, opt := range opts {
		opt(&e)
	}
	r.counters = append(r.counters, &e)
	return e.increment
}

// PublishAndReset publishes actual statistics and reset all internal state.
func (r *Registry) PublishAndReset(publisher Publisher) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.counters {
		v.publishAndReset(publisher)
	}
}

// WithLimit limits the number of the counters stored in the memory.
func WithLimit(limit int) func(counter *counter) {
	return func(counter *counter) {
		counter.limit = limit
	}
}
