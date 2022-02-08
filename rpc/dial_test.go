// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
	"storj.io/drpc/drpcmigrate"
)

func TestDialerUnencrypted(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	d := NewDefaultPooledDialer(nil)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ctx.Check(lis.Close)

	conn, err := d.DialAddressUnencrypted(ctx, lis.Addr().String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestDialHostnameVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	certificatePEM, privateKeyPEM := createTestingCertificate(t, "localhost")

	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	require.NoError(t, err)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	// start a server with the certificate
	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	serverAddr := tcpListener.Addr().String()

	listenMux := drpcmigrate.NewListenMux(tcpListener, len(drpcmigrate.DRPCHeader))

	listenCtx, listenCancel := context.WithCancel(ctx)
	defer listenCancel()
	ctx.Go(func() error {
		return listenMux.Run(listenCtx)
	})

	drpcListener := tls.NewListener(listenMux.Route(drpcmigrate.DRPCHeader), serverTLSConfig)
	defer ctx.Check(drpcListener.Close)

	acceptConnectionSuccess := func() error {
		conn, err := drpcListener.Accept()
		if err != nil {
			return errs.Wrap(err)
		}
		defer func() { _ = conn.Close() }()

		buffer := make([]byte, 256)
		_, err = conn.Read(buffer)
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			return errs.Wrap(err)
		}
		return nil
	}

	acceptConnectionFailure := func() error {
		conn, err := drpcListener.Accept()
		if err != nil {
			return errs.Wrap(err)
		}
		defer func() { _ = conn.Close() }()

		buffer := make([]byte, 256)
		_, err = conn.Read(buffer)
		if err == nil {
			return errs.New("expected connection failure, but there is no error")
		}
		return nil
	}

	// create client
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certificatePEM)

	// happy scenario 1 get hostname from address
	require.Empty(t, sync2.Concurrently(
		acceptConnectionSuccess,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs: certPool,
			}

			// use a domain name to ensure we can get hostname from the address
			localAddr := strings.ReplaceAll(serverAddr, "127.0.0.1", "localhost")
			conn, err := dialer.DialAddressHostnameVerification(ctx, localAddr)
			if err != nil {
				return errs.Wrap(err)
			}
			return errs.Wrap(conn.Close())
		},
	))

	// happy scenario 2
	require.Empty(t, sync2.Concurrently(
		acceptConnectionSuccess,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs:    certPool,
				ServerName: "localhost",
			}
			// Can't verify IPv6 during ci because of docker default
			// connection, err = dialer.DialAddressHostnameVerification(ctx, "[::1]:22111", clientTLSConfig)
			conn, err := dialer.DialAddressHostnameVerification(ctx, serverAddr)
			if err != nil {
				return errs.Wrap(err)
			}
			return errs.Wrap(conn.Close())
		},
	))

	// failure scenario invalid certificate
	require.Empty(t, sync2.Concurrently(
		acceptConnectionFailure,
		func() error {
			dialer := NewDefaultDialer(nil)
			dialer.HostnameTLSConfig = &tls.Config{
				RootCAs:    certPool,
				ServerName: "storj.test",
			}
			conn, err := dialer.DialAddressHostnameVerification(ctx, serverAddr)
			if err == nil {
				_ = conn.Close()
				return errs.New("expected an error")
			}
			if conn != nil {
				return errs.New("expected conn to be nil")
			}
			if !strings.Contains(err.Error(), "certificate is valid for localhost, not storj.test") {
				return errs.New("expected an error, got: %w", err)
			}
			return nil
		},
	))

	// test invalid hostname
	dialer := NewDefaultDialer(nil)
	_, err = dialer.DialAddressHostnameVerification(ctx, "storj.test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing port in address")
}

func createTestingCertificate(t *testing.T, hostname string) (certificatePEM []byte, privateKeyPEM []byte) {
	notAfter := time.Now().Add(1 * time.Minute)

	// first create a server certificate
	template := x509.Certificate{
		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames:              []string{hostname},
		SerialNumber:          big.NewInt(1337),
		BasicConstraintsValid: false,
		IsCA:                  true,
		NotAfter:              notAfter,
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	certificateDERBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	require.NoError(t, err)

	certificatePEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateDERBytes})

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	require.NoError(t, err)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})

	return certificatePEM, privateKeyPEM
}
