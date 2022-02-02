// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package telemetry

import (
	"context"
	"sync"
	"time"

	"storj.io/common/sync2"
)

// Reporter calls a function to report metrics periodically.
type Reporter struct {
	interval time.Duration
	mu       sync.Mutex
	cancel   context.CancelFunc
	stopped  bool

	send func(ctx context.Context) error
}

// NewReporter creates a reporter which calls send function once in each interval.
func NewReporter(interval time.Duration, send func(ctx context.Context) error) (rv *Reporter, err error) {
	return &Reporter{
		interval: interval,
		send:     send,
	}, nil
}

// Run calls Report roughly every Interval.
func (c *Reporter) Run(ctx context.Context) {
	c.mu.Lock()
	if c.stopped {
		c.mu.Unlock()
		return
	}
	ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	for {
		sync2.Sleep(ctx, jitter(c.interval))
		if ctx.Err() != nil {
			return
		}

		_ = c.Publish(ctx)
	}
}

// Publish bundles up all the current stats and writes them out as UDP packets.
func (c *Reporter) Publish(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)
	return c.send(ctx)
}
