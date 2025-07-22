// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information

// combiner package implements a combiner queue that allows pushing
// job items to a background worker that combines multiple job items
// into a single batch.
package combiner

import (
	"context"
	"sync"

	"storj.io/common/sync2"
)

// ProcessFunc processes a queue of jobs.
type ProcessFunc[Job any] func(ctx context.Context, queue *Queue[Job])

// Combiner combines multiple jobs for a single worker to process in batches.
type Combiner[Job any] struct {
	// ctx context to pass down to the handler.
	ctx    context.Context
	cancel context.CancelFunc

	// process is used for processing the jobs.
	process ProcessFunc[Job]
	fail    ProcessFunc[Job]
	// queueSize is size for the queues.
	queueSize int
	// workers contains all worker goroutines.
	workers sync2.WorkGroup

	// mu worker
	mu     sync.Mutex
	worker *worker[Job]
}

// Options is for configuring the combiner queue.
type Options[Job any] struct {
	Process   ProcessFunc[Job]
	Fail      ProcessFunc[Job]
	QueueSize int
}

// New creates a new combiner queue.
//
// Parent context is passed to the job processing as the context.
// process is used to process work items and fail is called for jobs
// when they need to be aborted, due to context cancellation.
// queueSize is used to create new queues, queueSize = -1 means the
// queue is unbounded.
func New[Job any](
	parent context.Context,
	options Options[Job],
) *Combiner[Job] {
	ctx, cancel := context.WithCancel(parent)
	c := &Combiner[Job]{
		ctx:       ctx,
		cancel:    cancel,
		process:   options.Process,
		fail:      options.Fail,
		queueSize: options.QueueSize,
	}
	return c
}

// Wait waits for the active workers to be completed.
func (combiner *Combiner[Job]) Wait(ctx context.Context) error {
	// TODO: this should be sensitive to context cancellation.
	combiner.workers.Wait()
	return nil
}

// Stop prevents new worker from being started, without
// canceling existing jobs.
func (combiner *Combiner[Job]) Stop() {
	combiner.workers.Close()
}

// Close shuts down all workers.
func (combiner *Combiner[Job]) Close() {
	combiner.cancel()
	combiner.Stop()
}

// Enqueue adds a new job to the queue.
func (combiner *Combiner[Job]) Enqueue(ctx context.Context, job Job) {
	combiner.mu.Lock()
	defer combiner.mu.Unlock()

	last := combiner.worker

	// Check whether we can use the last worker.
	if last != nil && last.jobs.TryPush(job) {
		// We've successfully added a job to an existing worker.
		return
	}

	// Create a new worker when one doesn't exist or the last one was full.
	next := &worker[Job]{
		jobs: NewQueue[Job](combiner.queueSize),
		done: make(chan struct{}),
	}
	combiner.worker = next
	if !next.jobs.TryPush(job) {
		// This should never happen.
		panic("invalid queue implementation")
	}

	// Start the worker.
	next.start(combiner)
}

// worker handles a batch of jobs.
type worker[Job any] struct {
	// jobs is an active queue of work tems to be completed.
	jobs *Queue[Job]
	// done is a channel that will be closed when the worker finishes.
	done chan struct{}
}

// schedule starts the worker.
func (worker *worker[Job]) start(combiner *Combiner[Job]) {
	// Try to add to worker pool, this may fail when we are shutting things down.
	workerStarted := combiner.workers.Go(func() {
		defer close(worker.done)
		// Ensure we fail any jobs that the handler didn't handle.
		defer func() {
			if !worker.jobs.Completed() {
				combiner.fail(combiner.ctx, worker.jobs)
			}
		}()

		// Handle the job queue.
		combiner.process(combiner.ctx, worker.jobs)
	})

	// If we failed to start a worker, then mark all the jobs as failures.
	if !workerStarted {
		combiner.fail(combiner.ctx, worker.jobs)
		close(worker.done)
	}
}
