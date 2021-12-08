// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.16
// +build go1.16

package quic

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/zeebo/errs"

	"storj.io/common/peertls/tlsopts"
)

const defaultIdleTimeout = 60 * time.Second

// Listener implements listener for QUIC.
type Listener struct {
	listener quic.Listener
	conn     *net.UDPConn
}

// NewListener returns a new listener instance for QUIC.
// The quic.Config may be nil, in that case the default values will be used.
// if the provided context is closed, all existing or following Accept calls will return an error.
func NewListener(conn *net.UDPConn, tlsConfig *tls.Config, quicConfig *quic.Config) (net.Listener, error) {
	if conn == nil {
		return nil, Error.New("underlying udp connection can't be nil")
	}
	if tlsConfig == nil {
		return nil, Error.New("tls config is not set")
	}
	tlsConfigCopy := tlsConfig.Clone()
	tlsConfigCopy.NextProtos = []string{tlsopts.StorjApplicationProtocol}

	if quicConfig == nil {
		quicConfig = &quic.Config{
			MaxIdleTimeout: defaultIdleTimeout,
			// disable address validation in QUIC (it costs an extra round-trip, and we believe
			// it to be unnecessary given the low potential for traffic amplification attacks).
			AcceptToken: func(clientAddr net.Addr, token *quic.Token) bool {
				return true
			},
		}
	}

	listener, err := quic.Listen(conn, tlsConfigCopy, quicConfig)
	if err != nil {
		return nil, err
	}

	return &Listener{
		listener: listener,
		conn:     conn,
	}, nil
}

// Accept waits for and returns the next available quic session to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	ctx := context.Background()
	session, err := l.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return &Conn{
		session: session,
	}, nil
}

// Close closes the QUIC listener.
func (l *Listener) Close() (err error) {
	return errs.Combine(l.listener.Close(), l.conn.Close())
}

// Addr returns the local network addr that the server is listening on.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}
