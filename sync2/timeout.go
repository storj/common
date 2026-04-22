// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2

import (
	"sync/atomic"
	"time"
)

// WithTimeout calls `do` and when the timeout is reached before `do`
// returns, it'll call `onTimeout` concurrently.
//
// If `do` returns at roughly the same instant the timer fires, exactly one
// of them wins: either `do`'s completion stops the timer and `onTimeout`
// is not called, or `onTimeout` runs and `do`'s subsequent return does not
// suppress it. This prevents a successful `do` from being followed by a
// spurious `onTimeout` call that can poison shared state (e.g. cancel a
// context that downstream operations still depend on).
//
// When WithTimeout returns it's guaranteed to not call onTimeout.
func WithTimeout(timeout time.Duration, do, onTimeout func()) {
	// state transitions: 0 (pending) -> 1 (do completed) OR 0 -> 2 (timed out).
	var state atomic.Int32
	done := make(chan struct{})
	t := time.AfterFunc(timeout, func() {
		defer close(done)
		if state.CompareAndSwap(0, 2) {
			onTimeout()
		}
	})
	defer func() {
		if state.CompareAndSwap(0, 1) {
			t.Stop()
			return
		}
		<-done
	}()
	do()
}
