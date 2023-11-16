// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpcpool

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
	"storj.io/drpc"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	calls := 0
	dial := func(ctx context.Context) (RawConn, *tls.ConnectionState, error) {
		calls++
		return &emptyConn{}, nil, nil
	}

	check := func(t *testing.T, pool *Pool, counts ...int) {
		calls = 0

		_, err := pool.Get(ctx, "key1", nil, dial)
		require.NoError(t, err)
		require.Equal(t, counts[0], calls)

		c1, err := pool.Get(ctx, "key1", nil, dial)
		require.NoError(t, err)
		require.Equal(t, counts[1], calls)

		c2, err := pool.Get(ctx, "key2", nil, dial)
		require.NoError(t, err)
		require.Equal(t, counts[2], calls)

		_ = c1.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, counts[3], calls)

		_ = c2.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, counts[4], calls)

		_ = c1.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, counts[5], calls)

		_, err = pool.Get(ctx, "key3", nil, dial)
		require.NoError(t, err)
		require.Equal(t, counts[6], calls)

		_, err = pool.Get(WithForceDial(ctx), "key4", nil, dial)
		require.NoError(t, err)
		require.Equal(t, counts[7], calls)
	}

	t.Run("Cached", func(t *testing.T) {
		check(t, New(Options{}), 0, 0, 0, 1, 2, 2, 2, 3)
	})

	t.Run("Nil", func(t *testing.T) {
		check(t, (*Pool)(nil), 0, 0, 0, 1, 2, 2, 2, 3)
	})
}

func TestExpired(t *testing.T) {
	ctx := testcontext.New(t)

	pool := New(Options{
		MaxLifetime: 1 * time.Hour,
	})

	dialed := 0
	dial := func(ctx context.Context) (RawConn, *tls.ConnectionState, error) {
		dialed++
		return &emptyConn{}, nil, nil
	}

	// this will initialize the first key
	conn, err := pool.Get(ctx, "key1", nil, dial)
	require.NoError(t, err)
	require.Equal(t, 0, dialed)

	// need the first invoke to be opened
	err = conn.Invoke(ctx, "somerpc", nil, nil, nil)
	require.NoError(t, err)
	require.Equal(t, 1, dialed)

	// still we use the same instance
	err = conn.Invoke(ctx, "somerpc", nil, nil, nil)
	require.NoError(t, err)
	require.Equal(t, 1, dialed)

	require.NoError(t, conn.Close())

	// we get the reference, and save the reference to the poolValue
	ref, err := pool.get(ctx, conn.(*poolConn).pk, nil)
	require.NoError(t, err)
	pool.put(conn.(*poolConn).pk, ref)
	require.NoError(t, err)
	require.Equal(t, 1, dialed)

	// connection is not yet closed
	require.False(t, ref.conn.(*emptyConn).closed)

	// here we make the connection expired.
	ref.created = time.Now().Add(-2 * time.Hour)

	// this will create a new connection (after DoInvoke) as existing is expired
	conn, err = pool.Get(ctx, "key1", nil, dial)
	require.NoError(t, err)
	err = conn.Invoke(ctx, "somerpc", nil, nil, nil)
	require.NoError(t, err)

	require.Equal(t, 2, dialed)
	require.True(t, ref.conn.(*emptyConn).closed)
}

// fakes for the test

type emptyConn struct {
	drpc.Conn
	closed bool
}

func (e *emptyConn) Close() error {
	e.closed = true
	return nil
}

func (*emptyConn) Closed() <-chan struct{} { return nil }
func (*emptyConn) Unblocked() <-chan struct{} {
	x := make(chan struct{})
	close(x)
	return x
}

func (*emptyConn) Transport() drpc.Transport {
	return nil
}

func (*emptyConn) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
	return nil
}
