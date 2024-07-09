// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"golang.org/x/sys/windows"
)

func osversion() (version int64, ok bool) {
	info := windows.RtlGetVersion()
	if info == nil {
		return 0, false
	}

	// Current maximum minor version is 3,
	// so the following computation should be fine.

	return int64(info.MajorVersion)*10 + int64(info.MinorVersion), true
}
