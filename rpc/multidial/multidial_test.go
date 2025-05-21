// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package multidial

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testPeer struct {
	t              testing.TB
	network        string
	address        string
	client, server net.Conn
}

func newTestPeer(t testing.TB, network, address string) *testPeer {
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	return &testPeer{
		t:       t,
		network: network,
		address: address,
		client:  client,
		server:  server,
	}
}

func (p *testPeer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	require.Equal(p.t, network, p.network)
	require.Equal(p.t, address, p.address)
	return p.client, nil
}

func TestBasic_ServerKeepsConnsOpen(t *testing.T) {
	testPeer1 := newTestPeer(t, "testnet", "testaddr")
	testPeer2 := newTestPeer(t, "testnet", "testaddr")
	conn, err := NewMultidialer(testPeer1.DialContext, testPeer2.DialContext).Dial("testnet", "testaddr")
	require.NoError(t, err)

	data := make([]byte, 10*1024)
	_, err = rand.Read(data)
	require.NoError(t, err)

	done := make(chan struct{}, 2)
	winningPeer := 0
	if data[0] >= 128 {
		winningPeer = 1
	}

	readsComplete := []chan struct{}{make(chan struct{}), make(chan struct{})}

	for peerIdx, peer := range []*testPeer{testPeer1, testPeer2} {
		go func(peerIdx int, peer *testPeer) {
			defer func() { done <- struct{}{} }()

			readbuf := make([]byte, 10*1024)
			_, err := io.ReadFull(peer.server, readbuf)
			require.NoError(t, err)
			require.Equal(t, readbuf, data)

			close(readsComplete[peerIdx])
			for _, ch := range readsComplete {
				<-ch
			}

			if peerIdx != winningPeer {
				return
			}

			_, err = peer.server.Write([]byte("response"))
			require.NoError(t, err)
			for range 10 {
				runtime.Gosched()
			}

			_, err = peer.server.Write([]byte("data!!!!"))
			require.NoError(t, err)
		}(peerIdx, peer)
	}

	for i := range 10 {
		_, err = conn.Write(data[i*1024 : (i+1)*1024])
		require.NoError(t, err)
	}

	buf := make([]byte, len("response"))
	_, err = io.ReadFull(conn, buf)
	require.NoError(t, err)
	require.Equal(t, string(buf), "response")

	_, err = io.ReadFull(conn, buf)
	require.NoError(t, err)
	require.Equal(t, string(buf), "data!!!!")

	require.NoError(t, conn.Close())

	<-done
	<-done
	require.NoError(t, testPeer1.server.Close())
	require.NoError(t, testPeer2.server.Close())
}

func TestBasic_ServerClosesUnneededConn(t *testing.T) {
	testPeer1 := newTestPeer(t, "testnet", "testaddr")
	testPeer2 := newTestPeer(t, "testnet", "testaddr")
	conn, err := NewMultidialer(testPeer1.DialContext, testPeer2.DialContext).Dial("testnet", "testaddr")
	require.NoError(t, err)

	data := make([]byte, 10*1024)
	_, err = rand.Read(data)
	require.NoError(t, err)

	done := make(chan struct{}, 2)
	winningPeer := 0
	if data[0] >= 128 {
		winningPeer = 1
	}
	closed := make(chan struct{})
	readsComplete := []chan struct{}{make(chan struct{}), make(chan struct{})}
	for peerIdx, peer := range []*testPeer{testPeer1, testPeer2} {
		go func(peerIdx int, peer *testPeer) {
			defer func() { done <- struct{}{} }()

			readbuf := make([]byte, 10*1024)
			_, err := io.ReadFull(peer.server, readbuf)
			require.NoError(t, err)
			require.Equal(t, readbuf, data)

			close(readsComplete[peerIdx])
			for _, ch := range readsComplete {
				<-ch
			}

			if peerIdx != winningPeer {
				require.NoError(t, peer.server.Close())
				close(closed)
				return
			}

			<-closed
			for range 10 {
				runtime.Gosched()
			}

			_, err = peer.server.Write([]byte("response"))
			require.NoError(t, err)

			_, err = peer.server.Write([]byte("data!!!!"))
			require.NoError(t, err)

			require.NoError(t, peer.server.Close())
		}(peerIdx, peer)
	}

	for i := range 10 {
		_, err = conn.Write(data[i*1024 : (i+1)*1024])
		require.NoError(t, err)
	}

	buf := make([]byte, len("response"))
	_, err = io.ReadFull(conn, buf)
	require.NoError(t, err)
	require.Equal(t, string(buf), "response")

	_, err = io.ReadFull(conn, buf)
	require.NoError(t, err)
	require.Equal(t, string(buf), "data!!!!")

	require.NoError(t, conn.Close())

	<-done
	<-done
}

func TestDeadlines(t *testing.T) {
	var conn1, conn2 mockConn
	conn1.written = make(chan struct{})
	conn, err := NewMultidialer(
		func(context.Context, string, string) (net.Conn, error) { return &conn1, nil },
		func(context.Context, string, string) (net.Conn, error) { return &conn2, nil },
	).Dial("testnet", "testaddr")
	require.NoError(t, err)

	conn1.sleep = make(chan struct{})

	t1 := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
	t2 := time.Date(2001, 1, 2, 3, 4, 5, 6, time.UTC)
	t3 := time.Date(2002, 1, 2, 3, 4, 5, 6, time.UTC)

	require.NoError(t, conn.SetDeadline(t1))
	require.NoError(t, conn.SetReadDeadline(t2))
	require.NoError(t, conn.SetWriteDeadline(t3))

	require.Equal(t, conn1.deadline, time.Time{})
	require.Equal(t, conn1.readDeadline, time.Time{})
	require.Equal(t, conn1.writeDeadline, time.Time{})
	require.Equal(t, conn2.deadline, t1)
	require.Equal(t, conn2.readDeadline, t2)
	require.Equal(t, conn2.writeDeadline, t3)

	// okay, let the calls to conn1 go through
	close(conn1.sleep)

	// flush out the calls to conn1
	_, err = conn.Write([]byte("hello"))
	require.NoError(t, err)
	<-conn1.written

	require.Equal(t, t1, conn1.deadline)
	require.Equal(t, t2, conn1.readDeadline)
	require.Equal(t, t3, conn1.writeDeadline)

	require.NoError(t, conn.Close())
}

type mockConn struct {
	net.Conn
	sleep                                 chan struct{}
	written                               chan struct{}
	deadline, readDeadline, writeDeadline time.Time
}

func (m *mockConn) SetDeadline(t time.Time) error {
	if m.sleep != nil {
		<-m.sleep
	}
	m.deadline = t
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	if m.sleep != nil {
		<-m.sleep
	}
	m.readDeadline = t
	return nil
}
func (m *mockConn) SetWriteDeadline(t time.Time) error {
	if m.sleep != nil {
		<-m.sleep
	}
	m.writeDeadline = t
	return nil
}

func (m *mockConn) Write(p []byte) (n int, err error) {
	if m.written != nil {
		close(m.written)
	}
	return len(p), nil
}

func (m *mockConn) Close() error { return nil }
