// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpcpool

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/drpc"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	calls := 0
	dial := func(ctx context.Context) (drpc.Conn, *tls.ConnectionState, error) {
		calls++
		return emptyConn{}, nil, nil
	}

	check := func(t *testing.T, pool *Pool, counts ...int) {
		calls = 0

		_, err := pool.Get(ctx, "key1", nil, dial)
		require.NoError(t, err)
		require.Equal(t, calls, counts[0])

		c1, err := pool.Get(ctx, "key1", nil, dial)
		require.NoError(t, err)
		require.Equal(t, calls, counts[1])

		c2, err := pool.Get(ctx, "key2", nil, dial)
		require.NoError(t, err)
		require.Equal(t, calls, counts[2])

		_ = c1.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, calls, counts[3])

		_ = c2.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, calls, counts[4])

		_ = c1.Invoke(ctx, "somerpc", nil, nil, nil)
		require.Equal(t, calls, counts[5])

		_, err = pool.Get(ctx, "key3", nil, dial)
		require.NoError(t, err)
		require.Equal(t, calls, counts[6])

		_, err = pool.Get(WithForceDial(ctx), "key4", nil, dial)
		require.NoError(t, err)
		require.Equal(t, calls, counts[7])
	}

	t.Run("Cached", func(t *testing.T) {
		check(t, New(Options{}), 0, 0, 0, 1, 2, 2, 2, 3)
	})

	t.Run("Nil", func(t *testing.T) {
		check(t, (*Pool)(nil), 0, 0, 0, 1, 2, 2, 2, 3)
	})
}

// fakes for the test

type emptyConn struct{ drpc.Conn }

func (emptyConn) Close() error            { return nil }
func (emptyConn) Closed() <-chan struct{} { return nil }

func (emptyConn) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
	return nil
}
