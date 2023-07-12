// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build windows

package quic

import (
	"errors"

	"golang.org/x/sys/windows"
)

func isMsgSizeErr(err error) bool {
	return errors.Is(err, windows.WSAEMSGSIZE)
}
