// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

// Package metrics implements a server which displays only read-only monitoring data.
package metrics

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/common/debug"
)

func init() {
	// zero out the http.DefaultServeMux net/http/pprof so unhelpfully
	// side-effected.
	*http.DefaultServeMux = http.ServeMux{}
}

// Config defines configuration for metrics server.
type Config struct {
	Address  string `internal:"true"`
	TLSKey   string
	TLSCert  string
	ClientCA string
}

// Server provides endpoints for debugging.
type Server struct {
	log *zap.Logger

	listener net.Listener
	server   http.Server
	mux      http.ServeMux

	*debug.PrometheusEndpoint
}

// NewListener configures a TLS listener based on configuration..
func NewListener(config Config) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(config.TLSCert, config.TLSKey)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	caCert, err := os.ReadFile(config.ClientCA)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cfg := &tls.Config{
		Certificates:     []tls.Certificate{cert},
		ClientCAs:        caCertPool,
		ClientAuth:       tls.RequireAndVerifyClientCert,
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	return tls.Listen("tcp", config.Address, cfg)

}

// NewServer returns a new debug.Server.
func NewServer(log *zap.Logger, listener net.Listener, registry *monkit.Registry, config Config) (*Server, error) {

	server := &Server{
		log:                log,
		listener:           listener,
		PrometheusEndpoint: debug.NewPrometheusEndpoint(registry),
	}

	server.server.Handler = &server.mux
	server.mux.HandleFunc("/metrics", server.PrometheusEndpoint.PrometheusMetrics)

	return server, nil
}

// Run starts the debug endpoint.
func (server *Server) Run(ctx context.Context) error {
	if server.listener == nil {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return debug.Error.Wrap(server.server.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()

		err := server.server.Serve(server.listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return debug.Error.Wrap(err)
	})
	return group.Wait()
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return debug.Error.Wrap(server.server.Close())
}
