// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package googleprofiler attaches google cloud profiler to private/process.
package googleprofiler

import (
	"flag"

	"cloud.google.com/go/profiler"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/common/process"
	"storj.io/common/version"
)

var (
	errorClass  = errs.Class("initializing profiler")
	serviceName = flag.String("debug.profilername", "", "provide the name of the peer to enable continuous cpu/mem profiling for")
)

func init() {
	process.SetProfiler(func(log *zap.Logger) error {
		name := *serviceName
		if name != "" {
			info := version.Build
			if err := profiler.Start(profiler.Config{
				Service:        name,
				ServiceVersion: info.Version.String(),
			}); err != nil {
				return errorClass.Wrap(err)
			}
			log.Debug("success debug profiler init")
		}
		return nil
	})
}
