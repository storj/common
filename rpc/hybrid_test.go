// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/memory"
	"storj.io/common/rpc"
	_ "storj.io/common/rpc/quic" // register quic connector
)

func TestSetTransferRate(t *testing.T) {
	hybridConnector := rpc.NewHybridConnector()
	hybridConnector.SetTransferRate(memory.GB)
	var n int
	for _, candidate := range hybridConnector.Connectors() {
		if connector, ok := candidate.Connector().(*rpc.TCPConnector); ok {
			assert.Equal(t, memory.GB, connector.TransferRate)
			n++
		}
	}
	assert.Greater(t, n, 0, "expected at least one *TCPConnector")
}

func TestSetSendDRPCMuxHeader(t *testing.T) {
	hybridConnector := rpc.NewHybridConnector()
	hybridConnector.SetSendDRPCMuxHeader(false)
	var n int
	for _, candidate := range hybridConnector.Connectors() {
		if connector, ok := candidate.Connector().(*rpc.TCPConnector); ok {
			assert.False(t, connector.SendDRPCMuxHeader)
			n++
		}
	}
	assert.Greater(t, n, 0, "expected at least one *TCPConnector")
}
