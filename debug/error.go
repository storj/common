// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package debug

import (
	"fmt"

	"github.com/spacemonkeygo/monkit/v3"

	"storj.io/drpc/drpcerr"
)

func init() {
	monkit.AddErrorNameHandler(func(err error) (string, bool) {
		if code := drpcerr.Code(err); code != 0 {
			return fmt.Sprintf("drpc_%d", code), true
		}
		return "", false
	})
}
