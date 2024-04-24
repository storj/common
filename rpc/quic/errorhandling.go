// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !windows

package quic

func isMsgSizeErr(err error) bool {
	// *nix doesn't return a size error from Accept.
	return false
}
