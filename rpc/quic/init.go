// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.16 && !noquic
// +build go1.16,!noquic

package quic

import (
	"storj.io/common/rpc"
)

const quicConnectorPriority = 20

func init() {
	rpc.RegisterCandidateConnectorType("quic", func() rpc.Connector {
		return NewDefaultConnector(nil)
	}, quicConnectorPriority)
}
