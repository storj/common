// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/zeebo/errs"

	"storj.io/common/memory"
	"storj.io/common/netutil"
	"storj.io/common/peertls/tlsopts"
)

// ConnectorConn is a type that creates a connection and establishes a tls
// session.
type ConnectorConn interface {
	net.Conn
	ConnectionState() tls.ConnectionState
}

// Connector is a type that creates a ConnectorConn, given an address and
// a tls configuration.
type Connector interface {
	// DialContext is called to establish a encrypted connection using tls.
	DialContext(ctx context.Context, tlsconfig *tls.Config, address string) (ConnectorConn, error)
}

// ConnectorAdapter represents a dialer that can establish a net.Conn.
type ConnectorAdapter struct {
	DialContext func(ctx context.Context, network, address string) (net.Conn, error)
}

// TCPConnector implements a dialer that creates an encrypted connection using tls.
type TCPConnector struct {
	// TCPUserTimeout controls what setting to use for the TCP_USER_TIMEOUT
	// socket option on dialed connections. Only valid on linux. Only set
	// if positive.
	TCPUserTimeout time.Duration

	// TransferRate limits all read/write operations to go slower than
	// the size per second if it is non-zero.
	TransferRate memory.Size

	dialer *ConnectorAdapter
}

// NewDefaultTCPConnector creates a new TCPConnector instance with provided tcp dialer.
// If no dialer is predefined, net.Dialer is used by default.
func NewDefaultTCPConnector(dialer *ConnectorAdapter) TCPConnector {
	if dialer == nil {
		dialer = &ConnectorAdapter{
			DialContext: new(net.Dialer).DialContext,
		}
	}

	return TCPConnector{
		TCPUserTimeout: 15 * time.Minute,
		dialer:         dialer,
	}
}

// DialContext creates a encrypted tcp connection using tls.
func (t TCPConnector) DialContext(ctx context.Context, tlsConfig *tls.Config, address string) (_ ConnectorConn, err error) {
	defer mon.Task()(&ctx)(&err)

	rawConn, err := t.DialContextUnencrypted(ctx, address)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	// perform the handshake racing with the context closing. we use a buffer
	// of size 1 so that the handshake can proceed even if no one is reading.
	errCh := make(chan error, 1)
	conn := tls.Client(rawConn, tlsConfig)
	go func() { errCh <- conn.Handshake() }()

	// see which wins and close the raw conn if there was any error. we can't
	// close the tls connection concurrently with handshakes or it sometimes
	// will panic. cool, huh?
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}
	if err != nil {
		_ = rawConn.Close()
		return nil, Error.Wrap(err)
	}

	return &tlsConnWrapper{
		Conn:       conn,
		underlying: rawConn,
	}, nil
}

// DialContextUnencrypted creates a raw tcp connection.
func (t TCPConnector) DialContextUnencrypted(ctx context.Context, address string) (_ net.Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	conn, err := t.dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	if tcpconn, ok := conn.(*net.TCPConn); t.TCPUserTimeout > 0 && ok {
		if err := netutil.SetUserTimeout(tcpconn, t.TCPUserTimeout); err != nil {
			return nil, errs.Combine(Error.Wrap(err), Error.Wrap(conn.Close()))
		}
	}

	return &timedConn{
		Conn: netutil.TrackClose(newDrpcHeaderConn(conn)),
		rate: t.TransferRate,
	}, nil
}

// QUICConnector implements a dialer that creates a quic connection.
type QUICConnector struct {
	transferRate memory.Size

	config *quic.Config
}

// NewDefaultQUICConnector instantiates a new instance of QUICConnector.
// If no quic configuration is provided, default value will be used.
func NewDefaultQUICConnector(quicConfig *quic.Config) QUICConnector {
	if quicConfig == nil {
		quicConfig = &quic.Config{
			MaxIdleTimeout: 15 * time.Minute,
		}
	}
	return QUICConnector{
		config: quicConfig,
	}
}

// DialContext creates a quic connection.
func (c QUICConnector) DialContext(ctx context.Context, tlsConfig *tls.Config, address string) (ConnectorConn, error) {
	if tlsConfig == nil {
		return nil, Error.New("tls config is not set")
	}
	tlsConfigCopy := tlsConfig.Clone()
	tlsConfigCopy.NextProtos = []string{tlsopts.StorjApplicationProtocol}

	sess, err := quic.DialAddrContext(ctx, address, tlsConfigCopy, c.config)
	if err != nil {
		return nil, err
	}

	stream, err := sess.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}

	quicConn := &QUICConn{
		session: sess,
		stream:  stream,
	}

	return &connectorConnWrapper{
		Conn: &timedConn{
			Conn: netutil.TrackClose(quicConn),
			rate: c.transferRate,
		},
		state: quicConn.ConnectionState(),
	}, nil
}

// SetTransferRate returns a QUIC connector with the given transfer rate.
func (c QUICConnector) SetTransferRate(rate memory.Size) QUICConnector {
	c.transferRate = rate
	return c
}

// TransferRate returns the transfer rate set on the connector.
func (c QUICConnector) TransferRate() memory.Size {
	return c.transferRate
}

// connectorConnWrapper is a wrapper around a net.Conn that has established a tls
// session.
// It converts a net.Conn to fulfill ConnectorConn interface.
type connectorConnWrapper struct {
	net.Conn
	state tls.ConnectionState
}

func (w *connectorConnWrapper) ConnectionState() tls.ConnectionState {
	return w.state
}
