// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package rpcpool

import (
	"context"
	"sync"

	"storj.io/drpc"
)

// poolConn is a wrapper around a drpc.Conn that we cache.
type poolConn struct {
	mu     sync.Mutex
	active int

	drpc.Conn
	pk   poolKey
	pool *Pool
}

// incActive increments the number of active RPCs on the connection.
func (c *poolConn) incActive() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.active++
}

// decActive decrements the number of RPCs on the connection.
func (c *poolConn) decActive() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.active--
}

// forceClose closes the underlying drpc.Conn.
func (c *poolConn) forceClose() error {
	return c.Conn.Close()
}

// Close checks to see if there are no active RPCs. If there are none, it places
// the connection into the pool for reuse. Otherwise, it closes the connection.
func (c *poolConn) Close() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.active != 0 {
		return c.forceClose()
	}

	c.pool.cache.Put(c.pk, c)
	return nil
}

// Invoke wraps drpc.Conn's Invoke method and keeps track of the number of
// active RPCs.
func (c *poolConn) Invoke(ctx context.Context, rpc string, in, out drpc.Message) (err error) {
	defer mon.Task()(&ctx)(&err)

	c.incActive()
	defer c.decActive()

	return c.Conn.Invoke(ctx, rpc, in, out)
}

// NewStream wraps drpc.Conn's NewStream method and keeps track of the number
// of active RPCs.
func (c *poolConn) NewStream(ctx context.Context, rpc string) (_ drpc.Stream, err error) {
	defer mon.Task()(&ctx)(&err)

	c.incActive()

	stream, err := c.Conn.NewStream(ctx, rpc)
	if err != nil {
		c.decActive()
		return nil, err
	}

	// the stream's done channel is closed when we're sure no reads/writes are
	// coming in for that stream anymore. it has been fully terminated.
	go func() {
		<-stream.Context().Done()
		c.decActive()
	}()

	return stream, nil
}
