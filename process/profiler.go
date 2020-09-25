// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"go.uber.org/zap"
)

var initProfiler func(log *zap.Logger) error

// SetProfiler sets the profiler for process package.
//
// It panics on multiple calls.
func SetProfiler(fn func(log *zap.Logger) error) {
	if initProfiler != nil {
		panic("profiler already set")
	}
	initProfiler = fn
}
