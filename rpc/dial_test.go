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
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/drpc/drpcmigrate"
)

func TestDialerUnencrypted(t *testing.T) {
	d := NewDefaultPooledDialer(nil)

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() { _ = lis.Close() }()

	conn, err := d.DialAddressUnencrypted(context.Background(), lis.Addr().String())
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func TestDialHostnameVerification(t *testing.T) {
	ctx := context.Background()

	certificatePEM, privateKeyPEM := createCertificate(t, "localhost")

	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	require.NoError(t, err)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	// start a server with the certificate
	tcpListener, err := net.Listen("tcp", ":22111")
	require.NoError(t, err)

	listenMux := drpcmigrate.NewListenMux(tcpListener, len(drpcmigrate.DRPCHeader))
	go func() {
		err := listenMux.Run(ctx)
		require.NoError(t, err)
	}()

	drpcListener := tls.NewListener(listenMux.Route(drpcmigrate.DRPCHeader), serverTLSConfig)
	defer func(drpcListener net.Listener) {
		err := drpcListener.Close()
		require.NoError(t, err)
	}(drpcListener)

	serverErrorChannel := make(chan error, 20)
	acceptConnectionSuccess := func() {
		connection, err := drpcListener.Accept()
		if err != nil {
			serverErrorChannel <- err
			return
		}

		buffer := make([]byte, 256)
		_, err = connection.Read(buffer)
		if err != nil {
			serverErrorChannel <- err
			return
		}
	}

	acceptConnectionFailure := func() {
		connection, err := drpcListener.Accept()
		if err != nil {
			serverErrorChannel <- err
			return
		}

		buffer := make([]byte, 256)
		_, err = connection.Read(buffer)
		if err == nil {
			serverErrorChannel <- errors.New("expected connection failure, but there is no error")
			return
		}
	}

	// create client
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certificatePEM)

	// happy scenario 1 get hostname from address
	dialer := NewDefaultPooledDialer(nil)
	dialer.HostnameTLSConfig = &tls.Config{
		RootCAs: certPool,
	}
	go acceptConnectionSuccess()
	connection, err := dialer.DialAddressHostnameVerification(ctx, "localhost:22111")
	require.NoError(t, err)
	require.NotNil(t, connection)

	// happy scenario 2
	dialer = NewDefaultPooledDialer(nil)
	dialer.HostnameTLSConfig = &tls.Config{
		RootCAs:    certPool,
		ServerName: "localhost",
	}
	go acceptConnectionSuccess()
	// Can't verify IPv6 during ci because of docker default
	// connection, err = dialer.DialAddressHostnameVerification(ctx, "[::1]:22111", clientTLSConfig)
	connection, err = dialer.DialAddressHostnameVerification(ctx, "127.0.0.2:22111")
	require.NoError(t, err)
	require.NotNil(t, connection)

	// failure scenario invalid certificate
	dialer = NewDefaultPooledDialer(nil)
	dialer.HostnameTLSConfig = &tls.Config{
		RootCAs:    certPool,
		ServerName: "example.com",
	}
	go acceptConnectionFailure()
	connection, err = dialer.DialAddressHostnameVerification(ctx, "127.0.0.1:22111")
	require.Error(t, err)
	require.EqualError(t, err, "rpc: x509: certificate is valid for localhost, not example.com")
	require.Nil(t, connection)

	require.Empty(t, serverErrorChannel, "Serverside errors occurred while testing the client")

	dialer = NewDefaultPooledDialer(nil)
	_, err = dialer.DialAddressHostnameVerification(ctx, "tweakers.net")
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing port in address")
}

func createCertificate(t *testing.T, hostname string) (certificatePEM []byte, privateKeyPEM []byte) {
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
