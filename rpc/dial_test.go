// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDialerUnencrypted(t *testing.T) {
	d := NewDefaultPooledDialer(nil)

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() { _ = lis.Close() }()

	conn, err := d.DialAddressUnencrypted(context.Background(), lis.Addr().String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}
