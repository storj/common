// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/lucas-clemente/quic-go"

	"storj.io/common/peertls/tlsopts"
)

// QUICListener implements listener for QUIC.
type QUICListener struct {
	listener quic.Listener
}

// NewQUICListener returns a new listener instance for QUIC.
// The quic.Config may be nil, in that case the default values will be used.
// if the provided context is closed, all existing or following Accept calls will return an error.
func NewQUICListener(tlsConfig *tls.Config, address string, quicConfig *quic.Config) (net.Listener, error) {
	if tlsConfig == nil {
		return nil, Error.New("tls config is not set")
	}
	tlsConfigCopy := tlsConfig.Clone()
	tlsConfigCopy.NextProtos = []string{tlsopts.StorjApplicationProtocol}

	listener, err := quic.ListenAddr(address, tlsConfigCopy, quicConfig)
	if err != nil {
		return nil, err
	}

	return &QUICListener{
		listener: listener,
	}, nil
}

// Accept waits for and returns the next available quic session to the listener.
func (l *QUICListener) Accept() (net.Conn, error) {
	ctx := context.Background()
	session, err := l.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return &QUICConn{
		session: session,
	}, nil
}

// Close closes the QUIC listener.
func (l *QUICListener) Close() error {
	return l.listener.Close()
}

// Addr returns the local network addr that the server is listening on.
func (l *QUICListener) Addr() net.Addr {
	return l.listener.Addr()
}
