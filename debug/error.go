// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"github.com/spacemonkeygo/monkit/v3"

	"storj.io/common/rpc/rpcstatus"
)

func init() {
	monkit.AddErrorNameHandler(func(err error) (string, bool) {
		var code uint64
	forLoop:
		for i := 0; i < 100; i++ {
			if v, ok := err.(interface{ Name() (string, bool) }); ok {
				if cls, ok := v.Name(); ok && cls != "" {
					return cls, true
				}
			}
			if v, ok := err.(interface{ Code() uint64 }); ok {
				if code == 0 {
					code = v.Code()
				}
			}
			switch v := err.(type) { //nolint: errorlint // this is a custom unwrap loop
			case interface{ Cause() error }:
				err = v.Cause()
			case interface{ Unwrap() error }:
				err = v.Unwrap()
			case interface{ Unwrap() []error }:
				errs := v.Unwrap()
				if len(errs) == 0 {
					break
				}
				err = errs[0]
			default:
				break forLoop
			}
		}
		if code != 0 {
			return "drpc_" + rpcstatus.StatusCode(code).String(), true
		}
		return "", false
	})
}
