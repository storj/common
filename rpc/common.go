// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"

	"storj.io/common/memory"
)

const (
	// IsDRPC is true if drpc is being used.
	IsDRPC = true
)

var mon = monkit.Package()

// Error wraps all of the errors returned by this package.
var Error = errs.Class("rpc")

//
// timed conns
//

// timedConn wraps a net.Conn so that all reads and writes get the specified timeout and
// return bytes no faster than the rate. If the timeout or rate are zero, they are
// ignored.
type timedConn struct {
	net.Conn
	rate memory.Size
}

// now returns time.Now if there's a nonzero rate.
func (t *timedConn) now() (now time.Time) {
	if t.rate > 0 {
		now = time.Now()
	}
	return now
}

// delay ensures that we sleep to keep the rate if it is nonzero. n is the number of
// bytes in the read or write operation we need to delay.
func (t *timedConn) delay(start time.Time, n int) {
	if t.rate > 0 {
		expected := time.Duration(n * int(time.Second) / t.rate.Int())
		if actual := time.Since(start); expected > actual {
			time.Sleep(expected - actual)
		}
	}
}

// Read wraps the connection read and adds sleeping to ensure the rate.
func (t *timedConn) Read(p []byte) (int, error) {
	start := t.now()
	n, err := t.Conn.Read(p)
	t.delay(start, n)
	return n, err
}

// Write wraps the connection write and adds sleeping to ensure the rate.
func (t *timedConn) Write(p []byte) (int, error) {
	start := t.now()
	n, err := t.Conn.Write(p)
	t.delay(start, n)
	return n, err
}

// NetConn returns the underlying conn, like *tls.Conn does.
func (t *timedConn) NetConn() net.Conn {
	return t.Conn
}

//
// tls conn wrapper
//

// tlsConnWrapper is a wrapper around a *tls.Conn that calls Close on the
// underlying connection when closed rather than trying to send a
// notification to the other side which may block forever.
type tlsConnWrapper struct {
	*tls.Conn
	underlying net.Conn
}

// Close closes the underlying connection.
func (t *tlsConnWrapper) Close() error { return t.underlying.Close() }

// NetConn returns the underlying conn, like *tls.Conn does.
func (t *tlsConnWrapper) NetConn() net.Conn {
	return t.underlying
}
