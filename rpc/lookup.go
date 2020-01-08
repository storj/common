// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"context"
	"net"
)

// LookupNodeAddress resolves a storage node address to the first IP address resolved.
// If an IP address is accidentally provided it is returned back. This function
// is used to resolve storage node IP addresses so that uplinks can use
// IP addresses directly without resolving many hosts.
func LookupNodeAddress(ctx context.Context, nodeAddress string) string {
	// We check if the address is an IP address.
	ip := net.ParseIP(nodeAddress)
	if ip == nil {
		// We have a hostname not an IP address so we should resolve the IP address
		// to give back to the uplink client.
		addresses, err := net.DefaultResolver.LookupHost(ctx, nodeAddress)
		if err == nil && len(addresses) > 0 {
			// We return the first address found because some DNS servers already do
			// round robin load balancing and we would be messing with their behaviour
			// if we tried to get smart here.
			return addresses[0]
		}
	}

	// We ignore the error because if this fails for some reason we can just
	// re-use the hostname, it just won't be as fast for the uplink to dial.
	return nodeAddress
}
