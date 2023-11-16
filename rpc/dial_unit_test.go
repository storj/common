// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/drpc"
	"storj.io/drpc/drpcconn"
)

func TestDialCloseIfError(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	connector := &trackedConnector{}

	d := Dialer{
		ConnectionOptions: drpcconn.Options{
			Manager: NewDefaultManagerOptions(),
		},
		Connector: connector,
	}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ctx.Check(lis.Close)

	conn, err := d.DialNode(ctx, storj.NodeURL{
		Address: lis.Addr().String(),
		NoiseInfo: storj.NoiseInfo{
			PublicKey: "asd",
			Proto:     storj.NoiseProto(-1),
		},
	}, DialOptions{
		ReplaySafe: true,
	})
	require.NoError(t, err)

	out := struct{}{}
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = conn.Invoke(tctx, "rpc", fakeEncoding{}, struct{}{}, &out)
	require.Error(t, err)

	require.True(t, connector.conn.closed)
}

type trackedConn struct {
	closed bool
}

func (t *trackedConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (t *trackedConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (t *trackedConn) Close() error {
	t.closed = true
	return nil
}

func (t *trackedConn) LocalAddr() net.Addr {
	return nil
}

func (t *trackedConn) RemoteAddr() net.Addr {
	return nil
}

func (t *trackedConn) SetDeadline(_ time.Time) error {
	return nil
}

func (t *trackedConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (t *trackedConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

func (t trackedConn) ConnectionState() tls.ConnectionState {
	return tls.ConnectionState{}
}

type trackedConnector struct {
	conn *trackedConn
}

func (t *trackedConnector) DialContextUnencrypted(context.Context, string) (net.Conn, error) {
	if t.conn != nil {
		return nil, errs.New("already connected")
	}
	t.conn = &trackedConn{}
	return t.conn, nil
}

func (t *trackedConnector) DialContextUnencryptedUnprefixed(context.Context, string) (net.Conn, error) {
	if t.conn != nil {
		return nil, errs.New("already connected")
	}
	t.conn = &trackedConn{}
	return t.conn, nil
}

func (t *trackedConnector) DialContext(ctx context.Context, tlsconfig *tls.Config, address string) (ConnectorConn, error) {
	if t.conn != nil {
		return nil, errs.New("already connected")
	}
	t.conn = &trackedConn{}
	return t.conn, nil
}

type fakeEncoding struct {
}

func (fakeEncoding) Marshal(msg drpc.Message) ([]byte, error) {
	return []byte{}, nil
}

func (fakeEncoding) Unmarshal(buf []byte, msg drpc.Message) error {
	return nil
}
