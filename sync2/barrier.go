// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2

import (
	"context"
	"sync/atomic"
)

// Barrier blocks callers of Wait until n of them have arrived.
//
// After that, all callers pass through without waiting, which means
// that calling Wait more than n times is safe.
type Barrier struct {
	remaining atomic.Int64
	fence     Fence
}

// NewBarrier creates a barrier for n callers.
func NewBarrier(n int) *Barrier {
	barrier := &Barrier{}
	barrier.remaining.Store(int64(n))
	if n <= 0 {
		barrier.fence.Release()
	}
	return barrier
}

// Wait marks the caller as arrived and waits until n callers have arrived.
// It returns the context error, when the context is canceled before that.
func (barrier *Barrier) Wait(ctx context.Context) error {
	if barrier.remaining.Add(-1) == 0 {
		barrier.fence.Release()
	}
	if !barrier.fence.Wait(ctx) {
		return ctx.Err()
	}
	return nil
}
