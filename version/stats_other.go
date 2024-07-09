// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !windows

package version

func osversion() (major, minor int64, ok bool) {
	return 0, 0, false
}
