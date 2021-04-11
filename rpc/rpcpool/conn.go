// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package rpcpool

import (
	"context"
	"sync"

	"github.com/zeebo/errs"

	"storj.io/drpc"
)

// poolConn is a wrapper around a drpc.Conn that we cache.
type poolConn struct {
	mu     sync.Mutex
	active int
	conn   drpc.Conn
	closed bool

	pk   poolKey
	pool *Pool

	dial Dialer
}

// forceClose closes the underlying drpc.Conn.
func (c *poolConn) forceClose() error {
	c.mu.Lock()
	c.closed = true
	conn := c.conn
	c.mu.Unlock()

	return conn.Close()
}

// Close checks to see if there are no active RPCs. If there are none, it places
// the connection into the pool for reuse. Otherwise, it closes the current
// live connection and prevents future ones from starting.
func (c *poolConn) Close() (err error) {
	c.mu.Lock()

	if c.active != 0 {
		c.closed = true
		conn := c.conn
		c.mu.Unlock()
		return conn.Close()
	}

	c.mu.Unlock()
	c.pool.cache.Put(c.pk, c)
	return nil
}

// Invoke wraps drpc.Conn's Invoke method and keeps track of the number of
// active RPCs, starting a new valid connection if necessary.
func (c *poolConn) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) (err error) {
	defer mon.Task()(&ctx)(&err)

	c.mu.Lock()

	// TODO: checking if the conn is open before the request is racy. it
	// would be better to act on the returned error of conn.Invoke to
	// find out if the connection closed before we started. We won't
	// be able to do anything if data started going over a Write, but
	// from an API perspective, drpc telling us the connection is closed
	// on an attempt cuts down on some possible avoidable races.

	conn, err := c.lockedGetConn(ctx)
	if err != nil {
		c.mu.Unlock()
		return err
	}

	c.active++
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		c.active--
		c.mu.Unlock()
	}()

	return conn.Invoke(ctx, rpc, enc, in, out)
}

// NewStream wraps drpc.Conn's NewStream method and keeps track of the number
// of active RPCs, creating a new connection if the current one is dead.
func (c *poolConn) NewStream(ctx context.Context, rpc string, enc drpc.Encoding) (_ drpc.Stream, err error) {
	defer mon.Task()(&ctx)(&err)

	c.mu.Lock()

	// TODO: checking if the conn is open before the request is racy. it
	// would be better to act on the returned error of conn.NewStream to
	// find out if the connection closed before we started. We won't
	// be able to do anything if data started going over a Write, but
	// from an API perspective, drpc telling us the connection is closed
	// on an attempt cuts down on some possible avoidable races.

	conn, err := c.lockedGetConn(ctx)
	if err != nil {
		c.mu.Unlock()
		return nil, err
	}

	c.active++
	c.mu.Unlock()

	stream, err := conn.NewStream(ctx, rpc, enc)
	if err != nil {
		c.mu.Lock()
		c.active--
		c.mu.Unlock()
		return nil, err
	}

	// the stream's done channel is closed when we're sure no reads/writes are
	// coming in for that stream anymore. it has been fully terminated.
	go func() {
		<-stream.Context().Done()
		c.mu.Lock()
		c.active--
		c.mu.Unlock()
	}()

	return stream, nil
}

func (c *poolConn) lockedGetConn(ctx context.Context) (drpc.Conn, error) {
	if c.conn.Closed() {
		if c.closed {
			return nil, errs.New("conn closed")
		}
		conn, err := c.dial(ctx)
		if err != nil {
			return nil, err
		}
		c.conn = conn
	}
	return c.conn, nil
}

// Closed returns if the conn is no longer usable.
func (c *poolConn) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// Stale returns if the conn will have to dial again to be used, or is
// no longer usable.
func (c *poolConn) Stale() bool {
	c.mu.Lock()
	conn := c.conn
	closed := c.closed
	c.mu.Unlock()
	return closed || conn.Closed()
}

// Transport returns the transport the conn is using.
func (c *poolConn) Transport() drpc.Transport {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	return conn.Transport() // okay if this is a closed one.
}
