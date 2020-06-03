// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package socket

import (
	"syscall"
)

// https://github.com/silviov/TCP-LEDBAT/ provides this, though there are other
// options. TODO: should we probe for available ones on startup?
const linuxLowPrioCongController = "ledbat"

func setLowPrioCongestionController(fd int) error {
	return syscall.SetsockoptString(fd, syscall.IPPROTO_TCP, syscall.TCP_CONGESTION, linuxLowPrioCongController)
}
