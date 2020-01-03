// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// Package rpctimeout provides helpers to have timeouts on rpc streams.
package rpctimeout

import (
	"context"
	"time"
)

// Run runs the provided function with a context that will be canceled after
// the provided duration or when the provided function returns. It returns either
// the error from the context or the error from the function. It runs the function
// in its own goroutine and DOES NOT wait for it to exit. This is on purpose to get
// around some grpc brain damage with respect to canceling operations on server streams.
func Run(ctx context.Context, timeout time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errch := make(chan error, 1)
	go func() { errch <- fn(ctx) }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errch:
		return err
	}
}
