// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/spacemonkeygo/monkit/v3"
	"go.uber.org/zap"

	"storj.io/private/debug"
)

var (
	// DebugAddrFlag for --debug.addr.
	DebugAddrFlag = flag.String("debug.addr", "127.0.0.1:0", "address to listen on for debug endpoints")
)

func initDebug(log *zap.Logger, r *monkit.Registry, atomicLevel *zap.AtomicLevel) (err error) {
	if *DebugAddrFlag == "" {
		return nil
	}

	ln, err := net.Listen("tcp", *DebugAddrFlag)
	if err != nil {
		return err
	}

	go func() {
		server := debug.NewServerWithAtomicLevel(log, ln, r, debug.Config{
			Address: *DebugAddrFlag,
		}, atomicLevel)
		log.Debug(fmt.Sprintf("debug server listening on %s", ln.Addr().String()))
		err := server.Run(context.TODO())
		if err != nil {
			log.Error("debug server died", zap.Error(err))
		}
	}()

	return nil
}
