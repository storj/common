// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package socket

import (
	"os"
	"syscall"
)

var linuxLowPrioCongController = os.Getenv("STORJ_SOCKET_LOWPRIO_CTL")

func setLowPrioCongestionController(fd int) error {
	if linuxLowPrioCongController != "" {
		return syscall.SetsockoptString(fd, syscall.IPPROTO_TCP, syscall.TCP_CONGESTION, linuxLowPrioCongController)
	}
	return nil
}

func setLowEffortQoS(fd int) error {
	return syscall.SetsockoptByte(fd, syscall.SOL_IP, syscall.IP_TOS, byte(dscpLE)<<2)
}
