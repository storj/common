// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

// Connectors returns the list of connectors for testing.
func (c *HybridConnector) Connectors() []candidateConnector { //nolint: revive // for testing
	return c.connectors
}

// Connector returns the actual connector for testing.
func (c candidateConnector) Connector() Connector {
	return c.connector
}
