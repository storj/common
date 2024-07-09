// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !windows

package version

func osversion() (version int64, ok bool) {
	return 0, false
}
