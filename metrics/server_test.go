// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package metrics

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"storj.io/common/testcontext"
)

func TestMetricsServer(t *testing.T) {
	ctx := testcontext.New(t)

	registry := monkit.NewRegistry()
	config := Config{
		Address:  ":1234",
		TLSKey:   "testdata/server.key",
		TLSCert:  "testdata/server.crt",
		ClientCA: "testdata/ca-client.crt",
	}

	listener, err := NewListener(config)
	require.NoError(t, err)
	fmt.Println(listener.Addr())
	srv, err := NewServer(zaptest.NewLogger(t), listener, registry, config)
	require.NoError(t, err)

	go func() {
		_ = srv.Run(ctx)
	}()

	registry.ScopeNamed("test").Chain(
		monkit.StatSourceFunc(func(cb func(key monkit.SeriesKey, field string, val float64)) {
			cb(monkit.NewSeriesKey("m1"), "f1", 1)
			cb(monkit.NewSeriesKey("m2"), "f2", 2)
			cb(monkit.NewSeriesKey("m1"), "f3", 3)
			cb(monkit.NewSeriesKey("m3"), "", 4)
		}))

	caCert, err := os.ReadFile("testdata/ca.crt")
	require.NoError(t, err)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t.Run("http client must fail", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET", "http://"+listener.Addr().String()+"/metrics", nil)
		require.NoError(t, err)

		client := &http.Client{}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()
		require.Equal(t, 400, resp.StatusCode)

	})

	t.Run("https client must fail without client ca", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/metrics", listener.Addr()), nil)
		require.NoError(t, err)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}

		resp, err := client.Do(req)
		require.ErrorContains(t, err, "failed to verify certificate")
		if resp != nil {
			_ = resp.Body.Close()
		}

	})

	t.Run("https client must work with client ca", func(t *testing.T) {
		parts := strings.Split(listener.Addr().String(), ":")
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://127.0.0.1:%s/metrics", parts[len(parts)-1]), nil)
		require.NoError(t, err)

		cert, err := tls.LoadX509KeyPair("testdata/client1.crt", "testdata/client1.key")
		require.NoError(t, err)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		}

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "gauge")

	})

}
