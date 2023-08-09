// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !go1.20 || noquic
// +build !go1.20 noquic

package quic

import (
	"crypto/tls"
	"net"
	"sync"
)

// Ignore unused warnings when building without quic.
var _ = isMsgSizeErr
var _ = mon

// Listener implements a stub/noop listener.
type Listener struct {
	isClosed    bool
	closedMutex sync.Mutex
	closedCond  *sync.Cond
}

// NewListener returns a new stub/noop listener. It will never return a connection.
func NewListener(conn *net.UDPConn, tlsConfig *tls.Config, quicConfig interface{}) (net.Listener, error) {
	l := &Listener{}
	l.closedCond = sync.NewCond(&l.closedMutex)
	return l, nil
}

// Accept simply blocks until the listener is closed.
func (l *Listener) Accept() (net.Conn, error) {
	l.closedMutex.Lock()
	defer l.closedMutex.Unlock()
	for !l.isClosed {
		l.closedCond.Wait()
	}
	return nil, net.ErrClosed
}

// Close closes the listener.
func (l *Listener) Close() (err error) {
	l.closedMutex.Lock()
	l.isClosed = true
	l.closedMutex.Unlock()
	l.closedCond.Broadcast()
	return nil
}

type dummyAddr struct{}

func (d dummyAddr) Network() string {
	return "dummy"
}

func (d dummyAddr) String() string {
	return "dummy"
}

// Addr returns the local network addr that the server is listening on.
func (l *Listener) Addr() net.Addr {
	return dummyAddr{}
}
