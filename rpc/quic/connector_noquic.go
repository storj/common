// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.21 || noquic
// +build go1.21 noquic

package quic

import (
	"context"
	"crypto/tls"

	"storj.io/common/memory"
	"storj.io/common/rpc"
)

// Connector implements a stub dialer that always fails.
type Connector struct{}

// NewDefaultConnector returns a stub connector that always fails.
func NewDefaultConnector(quicConfig interface{}) Connector {
	return Connector{}
}

// DialContext returns a failure.
func (c Connector) DialContext(ctx context.Context, tlsConfig *tls.Config, address string) (_ rpc.ConnectorConn, err error) {
	return nil, ErrQuicDisabled
}

// SetTransferRate has no effect.
func (c *Connector) SetTransferRate(rate memory.Size) {}

// TransferRate returns zero.
func (c Connector) TransferRate() memory.Size {
	return 0
}
