// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/testcontext"
)

func TestLookupNodeAddress_Host(t *testing.T) {
	// When we provide a host to LookupHostFirstAddress we should get a valid IP address back.
	address := LookupNodeAddress(context.Background(), "google.com")

	// Verify we get a properly formatted IP address back.
	ip := net.ParseIP(address)
	assert.NotNil(t, ip)
}

func TestLookupNodeAddress_HostAndPort(t *testing.T) {
	// When we provide a host to LookupHostFirstAddress we should get a valid IP address and port back.
	address := LookupNodeAddress(context.Background(), "google.com:8888")

	// Verify we get a properly formatted IP address back.
	host, port, err := net.SplitHostPort(address)
	assert.NoError(t, err)
	assert.Equal(t, "8888", port)
	assert.NotNil(t, net.ParseIP(host))
}

func TestLookupNodeAddress_IP(t *testing.T) {
	ctx := testcontext.New(t)

	tests := []string{
		"8.8.8.8",
		"2001:4860:4860::8888",
		"192.168.0.1:8888",
		"[2001:4860:4860::8888]:8888",
	}
	for _, test := range tests {
		address := LookupNodeAddress(ctx, test)
		assert.Equal(t, test, address)
	}
}
