// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information

// Package mpscqueue is a multi-producer, single-consumer queue.
package mpscqueue

import (
	"sync/atomic"
)

// Queue implements a lock-free multi-producer single-consumer queue.
//
// The design is based on http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue.
type Queue[T any] struct {
	head atomic.Pointer[node[T]]
	tail *node[T]
	stub node[T]

	noCopy noCopy //nolint:structcheck
}

type node[T any] struct {
	next  atomic.Pointer[node[T]]
	value T
}

// New creates a Queue.
func New[T any]() *Queue[T] {
	q := &Queue[T]{}
	q.Init()
	return q
}

// Init initializes the queue for receiving.
func (q *Queue[T]) Init() {
	q.head.Store(&q.stub)
	q.tail = &q.stub
}

// Enqueue adds a value to the queue.
//
// Enqueue is safe to call concurrently.
func (q *Queue[T]) Enqueue(value T) {
	n := &node[T]{value: value}
	old := q.head.Swap(n)
	old.next.Store(n)
}

// Dequeue receives a value from the queue, or returns ok=false if there weren't any values.
//
// Dequeue is not safe to call concurrently.
func (q *Queue[T]) Dequeue() (value T, ok bool) {
	next := q.tail.next.Load()
	if next == nil {
		return value, false
	}
	q.tail = next
	// zero the value to release whatever value it was to gc.
	value, next.value = next.value, value
	return value, true
}

// see sync2.noCopy for details.
type noCopy struct{}

func (noCopy) Lock() {}
