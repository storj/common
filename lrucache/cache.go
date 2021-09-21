// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package lrucache

import (
	"container/list"
	"sync"
	"time"
)

// Options controls the details of the expiration policy.
type Options struct {
	// Expiration is how long an entry will be valid. It is not
	// affected by LRU or anything: after this duration, the object
	// is invalidated. A non-positive value means no expiration.
	Expiration time.Duration

	// Capacity is how many objects to keep in memory.
	Capacity int
}

// cacheState contains all of the state for a cached entry.
type cacheState struct {
	once  sync.Once
	when  time.Time
	order *list.Element
	value interface{}
}

// ExpiringLRU caches values for string keys with a time based expiration and
// an LRU based eviciton policy.
type ExpiringLRU struct {
	mu    sync.Mutex
	opts  Options
	data  map[string]*cacheState
	order *list.List
}

// New constructs an ExpiringLRU with the given options.
func New(opts Options) *ExpiringLRU {
	return &ExpiringLRU{
		opts:  opts,
		data:  make(map[string]*cacheState, opts.Capacity),
		order: list.New(),
	}
}

// Get returns the value for some key if it exists and is valid. If not
// it will call the provided function. Concurrent calls will dedupe as
// best as they are able. If the function returns an error, it is not
// cached and further calls will try again.
func (e *ExpiringLRU) Get(key string, fn func() (interface{}, error)) (
	value interface{}, err error) {
	if e.opts.Capacity <= 0 {
		return fn()
	}

	for {
		e.mu.Lock()

		state, ok := e.data[key]
		switch {
		case !ok:
			for len(e.data) >= e.opts.Capacity {
				back := e.order.Back()
				delete(e.data, back.Value.(string))
				e.order.Remove(back)
			}
			state = &cacheState{
				when:  time.Now(),
				order: e.order.PushFront(key),
			}
			e.data[key] = state

		case e.opts.Expiration > 0 && time.Since(state.when) > e.opts.Expiration:
			delete(e.data, key)
			e.order.Remove(state.order)
			e.mu.Unlock()
			continue

		default:
			e.order.MoveToFront(state.order)
		}

		e.mu.Unlock()

		called := false
		state.once.Do(func() {
			called = true
			value, err = fn()

			if err == nil {
				// careful because we don't want a `(*T)(nil) != nil` situation
				// that's why we only assign to state.value if err == nil.
				state.value = value
			} else {
				// the once has been used. delete it so that any other waiters
				// will retry.
				e.mu.Lock()
				if e.data[key] == state {
					delete(e.data, key)
					e.order.Remove(state.order)
				}
				e.mu.Unlock()
			}
		})

		if called || state.value != nil {
			return state.value, err
		}
	}
}

// Delete explicitly removes a key from the cache if it exists.
func (e *ExpiringLRU) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	state, ok := e.data[key]
	if !ok {
		return
	}
	delete(e.data, key)
	e.order.Remove(state.order)
}

// Add adds a value to the cache.
//
// replaced is true if the key already existed in the cache and was valid, hence
// the value is replaced.
func (e *ExpiringLRU) Add(key string, value interface{}) (replaced bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	state, evicted := e.peek(key)
	if state == nil {
		if !evicted && e.order.Len() >= e.opts.Capacity {
			item := e.order.Back()
			delete(e.data, item.Value.(string))
		}

		e.data[key] = &cacheState{
			when:  time.Now(),
			order: e.order.PushFront(key),
			value: value,
		}

		return false
	}

	e.order.Remove(state.order)
	e.data[key] = &cacheState{
		when:  time.Now(),
		order: e.order.PushFront(key),
		value: value,
	}

	return true
}

// GetCached returns the value associated with key and true if it exists and
// hasn't expired, otherwise nil and false.
func (e *ExpiringLRU) GetCached(key string) (value interface{}, cached bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	state, _ := e.peek(key)
	if state == nil {
		return nil, false
	}

	e.order.MoveToFront(state.order)
	return state.value, true
}

// peek returns the state associated to the key if exists and it's valid,
// otherwise nil. evicted is true when the key existed but the state has
// expired.
//
// peek doesn't update the key as being recently used.
//
// NOTE the caller must always lock and unlock the mutex before calling this
// method.
func (e *ExpiringLRU) peek(key string) (state *cacheState, evicted bool) {
	state, ok := e.data[key]
	if !ok {
		return nil, false
	}

	if e.opts.Expiration > 0 && time.Since(state.when) > e.opts.Expiration {
		e.order.Remove(state.order)
		delete(e.data, key)

		return nil, true
	}

	return state, false
}
