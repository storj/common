// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupNodeAddressWithHost(t *testing.T) {
	// When we provide a host to LookupHostFirstAddress we should get a valid IP address back.
	address := LookupNodeAddress("google.com")

	// Verify we get a properly formatted IP address back.
	ip := net.ParseIP(address)
	assert.NotNil(t, ip)
}

func TestLookupNodeAddressWithIP(t *testing.T) {
	// When we provide an IP address to LookupHostFirstAddress we should get the same IP address back.
	address := LookupNodeAddress("8.8.8.8")

	// Verify we get the same IP address back.
	assert.Equal(t, "8.8.8.8", address)
}
