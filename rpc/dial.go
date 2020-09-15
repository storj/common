// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc/rpcpool"
	"storj.io/common/rpc/rpctracing"
	"storj.io/common/storj"
	"storj.io/drpc"
	"storj.io/drpc/drpcconn"
	"storj.io/drpc/drpcmanager"
	"storj.io/drpc/drpcstream"
)

// NewDefaultManagerOptions returns the default options we use for drpc managers.
func NewDefaultManagerOptions() drpcmanager.Options {
	return drpcmanager.Options{
		WriterBufferSize: 1024,
		Stream: drpcstream.Options{
			SplitSize: (4096 * 2) - 256,
		},
	}
}

// Dialer holds configuration for dialing.
type Dialer struct {
	// TLSOptions controls the tls options for dialing. If it is nil, only
	// insecure connections can be made.
	TLSOptions *tlsopts.Options

	// DialTimeout causes all the tcp dials to error if they take longer
	// than it if it is non-zero.
	DialTimeout time.Duration

	// DialLatency sleeps this amount if it is non-zero before every dial.
	// The timeout runs while the sleep is happening.
	DialLatency time.Duration

	// PoolOptions controls options for the connection pool.
	PoolOptions rpcpool.Options

	// ConnectionOptions controls the options that we pass to drpc connections.
	ConnectionOptions drpcconn.Options

	// Connector is how sockets are opened. If nil, net.Dialer is used.
	Connector Connector
}

// NewDefaultDialer returns a Dialer with default timeouts set.
func NewDefaultDialer(tlsOptions *tlsopts.Options) Dialer {
	return Dialer{
		TLSOptions:  tlsOptions,
		DialTimeout: 20 * time.Second,
		PoolOptions: rpcpool.Options{
			Capacity:       5,
			IdleExpiration: 2 * time.Minute,
		},
		ConnectionOptions: drpcconn.Options{
			Manager: NewDefaultManagerOptions(),
		},
		Connector: NewDefaultTCPConnector(nil),
	}
}

// DialNodeURL dials to the specified node url and asserts it has the given node id.
func (d Dialer) DialNodeURL(ctx context.Context, nodeURL storj.NodeURL) (_ *Conn, err error) {
	defer mon.Task()(&ctx, "node: "+nodeURL.ID.String()[0:8])(&err)

	if d.TLSOptions == nil {
		return nil, Error.New("tls options not set when required for this dial")
	}

	return d.dial(ctx, nodeURL.Address, d.TLSOptions.ClientTLSConfig(nodeURL.ID))
}

// DialAddressInsecure dials to the specified address and does not check the node id.
func (d Dialer) DialAddressInsecure(ctx context.Context, address string) (_ *Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	if d.TLSOptions == nil {
		return nil, Error.New("tls options not set when required for this dial")
	}

	return d.dial(ctx, address, d.TLSOptions.UnverifiedClientTLSConfig())
}

// DialAddressUnencrypted dials to the specified address without tls.
func (d Dialer) DialAddressUnencrypted(ctx context.Context, address string) (_ *Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	return d.dialUnencrypted(ctx, address)
}

// dial performs the dialing to the drpc endpoint with tls.
func (d Dialer) dial(ctx context.Context, address string, tlsConfig *tls.Config) (_ *Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	// include the timeout here so that it includes all aspects of the dial
	if d.DialTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, d.DialTimeout)
		defer cancel()
	}

	pool := rpcpool.New(d.PoolOptions, func(ctx context.Context) (drpc.Transport, error) {
		return d.dialTransport(ctx, address, tlsConfig)
	})

	conn, err := d.dialTransport(ctx, address, tlsConfig)
	if err != nil {
		return nil, err
	}
	state := conn.ConnectionState()

	if err := pool.Put(drpcconn.New(conn)); err != nil {
		return nil, err
	}

	return &Conn{
		state: state,
		Conn:  rpctracing.NewTracingWrapper(pool),
	}, nil
}

// dialTransport performs dialing to the drpc endpoint with tls.
func (d Dialer) dialTransport(ctx context.Context, address string, tlsConfig *tls.Config) (_ ConnectorConn, err error) {
	defer mon.Task()(&ctx)(&err)

	if d.DialLatency > 0 {
		timer := time.NewTimer(d.DialLatency)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return nil, Error.Wrap(ctx.Err())
		}
	}

	return d.Connector.DialContext(ctx, tlsConfig, address)
}

// dialUnencrypted performs dialing to the drpc endpoint with no tls.
func (d Dialer) dialUnencrypted(ctx context.Context, address string) (_ *Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	// include the timeout here so that it includes all aspects of the dial
	if d.DialTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, d.DialTimeout)
		defer cancel()
	}

	conn := rpcpool.New(d.PoolOptions, func(ctx context.Context) (drpc.Transport, error) {
		return d.dialTransportUnencrypted(ctx, address)
	})
	return &Conn{
		Conn: rpctracing.NewTracingWrapper(conn),
	}, nil
}

type unencryptedConnector interface {
	DialContextUnencrypted(context.Context, string) (net.Conn, error)
}

// dialTransportUnencrypted performs dialing to the drpc endpoint with no tls.
func (d Dialer) dialTransportUnencrypted(ctx context.Context, address string) (_ net.Conn, err error) {
	defer mon.Task()(&ctx)(&err)

	if unencConnector, ok := d.Connector.(unencryptedConnector); ok {
		// open the tcp socket to the address
		conn, err := unencConnector.DialContextUnencrypted(ctx, address)
		if err != nil {
			return nil, Error.Wrap(err)
		}

		return conn, nil
	}

	return nil, Error.New("unsupported transport type: %T, use TCPTransport", d.Connector)
}
