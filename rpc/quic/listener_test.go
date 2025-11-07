// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !noquic

package quic_test

import (
	"bytes"
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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/identity"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc/quic"
	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestListener(t *testing.T) {
	ctx := testcontext.New(t)

	certificatePEM, privateKeyPEM := createTestingCertificate(t, "localhost")

	certificate, err := tls.X509KeyPair(certificatePEM, privateKeyPEM)
	require.NoError(t, err)

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	require.NoError(t, err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	require.NoError(t, err)
	serverAddr = serverConn.LocalAddr().(*net.UDPAddr)

	listener, err := quic.NewListener(serverConn, serverTLSConfig, nil)
	require.NoError(t, err)
	require.NotNil(t, listener)

	clientData := []byte{4, 5, 6}
	serverData := []byte{1, 2, 3}

	errs := sync2.Concurrently(
		func() error {
			conn, err := listener.Accept()
			if err != nil {
				return errs.New("server Accept: %w", err)
			}

			fromClient := make([]byte, 3)
			n, err := conn.Read(fromClient)
			fromClient = fromClient[:n]
			if err != nil {
				return errs.New("server Read: %w", err)
			}
			if !bytes.Equal(fromClient, clientData) {
				return errs.New("server %v != %v", fromClient, clientData)
			}

			_, err = conn.Write(serverData)
			if err != nil {
				return errs.New("server Write: %w", err)
			}

			//HACKFIX: if we call Close before the client has received the answer
			// we'll end up discarding the buffer
			_, err = conn.Read([]byte{0})
			if errors.Is(err, io.EOF) {
				err = nil
			}
			if err != nil {
				return errs.New("server Wait: %w", err)
			}

			return errs.Wrap(conn.Close())
		}, func() error {
			ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{
				Difficulty:  0,
				Concurrency: 1,
			})
			if err != nil {
				return errs.New("could not generate an identity: %v", err)
			}
			tlsOptions, err := tlsopts.NewOptions(ident, tlsopts.Config{
				PeerIDVersions: "*",
			}, nil)
			if err != nil {
				return errs.New("could not get tls options: %v", err)
			}
			unverifiedClientConfig := tlsOptions.UnverifiedClientTLSConfig()

			connector := quic.NewDefaultConnector(nil)
			conn, err := connector.DialContext(ctx, unverifiedClientConfig, serverAddr.String())
			if err != nil {
				return errs.New("client DialContext: %w", err)
			}
			_, err = conn.Write(clientData)
			if err != nil {
				return errs.New("client Write: %w", err)
			}

			fromServer := make([]byte, 3)
			n, err := conn.Read(fromServer)
			fromServer = fromServer[:n]
			if errors.Is(err, io.EOF) && n > 0 {
				err = nil
			}
			if err != nil {
				return errs.New("client Read: %w", err)
			}

			if !bytes.Equal(fromServer, serverData) {
				return errs.New("client %v != %v", fromServer, serverData)
			}

			return errs.Wrap(conn.Close())
		})
	for _, err := range errs {
		require.NoError(t, err)
	}
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
