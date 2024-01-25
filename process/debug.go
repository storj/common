// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"context"
	"fmt"
	"net"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"storj.io/common/cfgstruct"
	"storj.io/common/debug"
)

var debugConfig struct {
	Debug debug.Config
}

func init() {
	cfgstruct.Bind(pflag.CommandLine, &debugConfig)
}

func initDebug(log *zap.Logger, r *monkit.Registry, atomicLevel *zap.AtomicLevel) (err error) {
	if debugConfig.Debug.Addr == "" {
		return nil
	}

	ln, err := net.Listen("tcp", debugConfig.Debug.Addr)
	if err != nil {
		return err
	}

	go func() {
		server := debug.NewServerWithAtomicLevel(log, ln, r, debugConfig.Debug, atomicLevel)
		log.Debug(fmt.Sprintf("debug server listening on %s", ln.Addr().String()))
		err := server.Run(context.TODO())
		if err != nil {
			log.Error("debug server died", zap.Error(err))
		}
	}()

	return nil
}
