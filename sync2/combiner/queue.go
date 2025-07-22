// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information

package combiner

import (
	"iter"
	"sync"
)

// Queue is a finalizable list of jobs with a limit to how many jobs it can handle.
type Queue[Job any] struct {
	// maxJobsPerBatch determines how many jobs can be pushed until the
	// queue is automatically closed.
	//
	// maxJobsPerBatch < 0 means there is no limit.
	maxJobsPerBatch int

	mu sync.Mutex
	// done indicates that no more items will be appended to the queue.
	done bool
	// list contains uncompleted jobs.
	list []Job
}

// NewQueue returns a new limited job queue.
func NewQueue[Job any](maxJobsPerBatch int) *Queue[Job] {
	return &Queue[Job]{
		maxJobsPerBatch: maxJobsPerBatch,
	}
}

// Completed returns true when the queue does not accept any new
// jobs and all the jobs have been completed.
func (jobs *Queue[Job]) Completed() bool {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()
	return jobs.done && len(jobs.list) == 0
}

// TryPush tries to add a job to the queue and returns
// false when the queue does not accept new jobs.
//
// maxJobsPerBatch < 0, means no limit.
func (jobs *Queue[Job]) TryPush(job Job) bool {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()

	// check whether we have finished work with this jobs queue.
	if jobs.done {
		return false
	}

	// check whether the queue is at capacity
	if jobs.maxJobsPerBatch >= 0 && len(jobs.list)+1 >= jobs.maxJobsPerBatch {
		jobs.done = true
	}

	jobs.list = append(jobs.list, job)
	return true
}

// PopAll returns all the jobs in this list.
//
// When there's no more items to be pulled, the queue automatically closes.
func (jobs *Queue[Job]) PopAll() (_ []Job, ok bool) {
	jobs.mu.Lock()
	defer jobs.mu.Unlock()

	// when we try to pop and the queue is empty, make the queue final.
	if len(jobs.list) == 0 {
		jobs.done = true
		return nil, false
	}

	list := jobs.list
	jobs.list = nil
	return list, true
}

// Batches iterates over batches until all done.
//
// The iterator slice should not be used outside of the loop.
func (jobs *Queue[Job]) Batches() iter.Seq[[]Job] {
	return func(yield func([]Job) bool) {
		for {
			all, ok := jobs.PopAll()
			if !ok {
				return
			}
			yield(all)
		}
	}
}
