// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !go1.20 || noquic

package quic

import (
	"crypto/tls"
	"net"
	"time"

	"storj.io/common/rpc"
)

// IsSupported returns whether quic building is enabled.
const IsSupported = false

// Conn is a stub/noop connection object.
type Conn struct{}

// Read panics.
func (c *Conn) Read(b []byte) (n int, err error) {
	panic("quic is disabled. how did you get here?")
}

// Write panics.
func (c *Conn) Write(b []byte) (_ int, err error) {
	panic("quic is disabled. how did you get here?")
}

// ConnectionState panics.
func (c *Conn) ConnectionState() tls.ConnectionState {
	panic("quic is disabled. how did you get here?")
}

// Close panics.
func (c *Conn) Close() error {
	panic("quic is disabled. how did you get here?")
}

// LocalAddr panics.
func (c *Conn) LocalAddr() net.Addr {
	panic("quic is disabled. how did you get here?")
}

// RemoteAddr panics.
func (c *Conn) RemoteAddr() net.Addr {
	panic("quic is disabled. how did you get here?")
}

// SetReadDeadline panics.
func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("quic is disabled. how did you get here?")
}

// SetWriteDeadline panics.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("quic is disabled. how did you get here?")
}

// SetDeadline panics.
func (c *Conn) SetDeadline(t time.Time) error {
	panic("quic is disabled. how did you get here?")
}

// TrackClose has no effect.
func TrackClose(conn rpc.ConnectorConn) rpc.ConnectorConn {
	return conn
}
