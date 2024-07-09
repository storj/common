// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"golang.org/x/sys/windows"
)

func osversion() (major, minor int64, ok bool) {
	info := windows.RtlGetVersion()
	if info == nil {
		return 0, 0, false
	}
	return int64(info.MajorVersion), int64(info.MinorVersion), true
}
