// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"

	"storj.io/common/identity"
	"storj.io/drpc"
)

// Conn is a wrapper around a drpc client connection.
type Conn struct {
	state tls.ConnectionState
	drpc.Conn
}

// Close closes the connection.
func (c *Conn) Close() error { return c.Conn.Close() }

// ConnectionState returns the tls connection state.
func (c *Conn) ConnectionState() tls.ConnectionState { return c.state }

// PeerIdentity returns the peer identity on the other end of the connection.
func (c *Conn) PeerIdentity() (*identity.PeerIdentity, error) {
	return identity.PeerIdentityFromChain(c.state.PeerCertificates)
}

// QUICConn is a wrapper around a quic connection and fulfills net.Conn interface.
type QUICConn struct {
	once sync.Once
	// The QUICConn.stream varible should never be directly accessed.
	// Always use QUICConn.getStream() instead.
	stream quic.Stream

	acceptErr error
	session   quic.Session
}

// Read implements the Conn Read method.
func (c *QUICConn) Read(b []byte) (n int, err error) {
	stream, err := c.getStream()
	if err != nil {
		return 0, err
	}
	return stream.Read(b)
}

// Write implements the Conn Write method.
func (c *QUICConn) Write(b []byte) (int, error) {
	stream, err := c.getStream()
	if err != nil {
		return 0, err
	}
	return stream.Write(b)
}

func (c *QUICConn) getStream() (quic.Stream, error) {
	if c.stream == nil {
		// When this function completes, it guarantees either c.err is not nil or c.stream is not nil
		c.once.Do(func() {
			stream, err := c.session.AcceptStream(context.Background())
			if err != nil {
				c.acceptErr = err
				return
			}

			c.stream = stream
		})
		if c.acceptErr != nil {
			return nil, c.acceptErr
		}
	}

	return c.stream, nil
}

// ConnectionState converts quic session state to tls connection state and returns tls state.
func (c *QUICConn) ConnectionState() tls.ConnectionState {
	state := c.session.ConnectionState()
	return tls.ConnectionState{
		Version:                     state.Version,
		HandshakeComplete:           state.HandshakeComplete,
		DidResume:                   state.DidResume,
		CipherSuite:                 state.CipherSuite,
		NegotiatedProtocol:          state.NegotiatedProtocol,
		ServerName:                  state.ServerName,
		PeerCertificates:            state.PeerCertificates,
		VerifiedChains:              state.VerifiedChains,
		SignedCertificateTimestamps: state.SignedCertificateTimestamps,
		OCSPResponse:                state.OCSPResponse,
	}
}

// Close closes the quic connection.
func (c *QUICConn) Close() error {
	return c.session.CloseWithError(quic.ErrorCode(0), "")
}

// LocalAddr returns the local address.
func (c *QUICConn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

// RemoteAddr returns the address of the peer.
func (c *QUICConn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
func (c *QUICConn) SetReadDeadline(t time.Time) error {
	stream, err := c.getStream()
	if err != nil {
		return err
	}
	return stream.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
func (c *QUICConn) SetWriteDeadline(t time.Time) error {
	stream, err := c.getStream()
	if err != nil {
		return err
	}
	return stream.SetWriteDeadline(t)
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
func (c *QUICConn) SetDeadline(t time.Time) error {
	stream, err := c.getStream()
	if err != nil {
		return err
	}

	return stream.SetDeadline(t)
}
