// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"time"

	"storj.io/drpc"
)

// ConnectionWithLatency wraps the original connection and add certain latency to it.
func ConnectionWithLatency(conn drpc.Conn, duration time.Duration) drpc.Conn {
	return &MessageInterceptor{
		delegate: conn,
		ResponseHook: func(rpc string, message drpc.Message, err error) {
			time.Sleep(duration)
		},
	}
}
