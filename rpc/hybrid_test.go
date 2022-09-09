// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/memory"
)

func TestSetTransferRate(t *testing.T) {
	hybridConnector := NewHybridConnector()
	hybridConnector.SetTransferRate(memory.GB)
	var n int
	for _, candidate := range hybridConnector.connectors {
		if connector, ok := candidate.connector.(*TCPConnector); ok {
			assert.Equal(t, memory.GB, connector.TransferRate)
			n++
		}
	}
	assert.Greater(t, n, 0, "expected at least one *TCPConnector")
}

func TestSetSendDRPCMuxHeader(t *testing.T) {
	hybridConnector := NewHybridConnector()
	hybridConnector.SetSendDRPCMuxHeader(false)
	var n int
	for _, candidate := range hybridConnector.connectors {
		if connector, ok := candidate.connector.(*TCPConnector); ok {
			assert.False(t, connector.SendDRPCMuxHeader)
			n++
		}
	}
	assert.Greater(t, n, 0, "expected at least one *TCPConnector")
}
