// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// +build linux

package netutil

import (
	"net"
	"time"

	"golang.org/x/sys/unix"
)

// SetUserTimeout sets the TCP_USER_TIMEOUT setting on the provided conn.
func SetUserTimeout(conn *net.TCPConn, timeout time.Duration) error {
	// By default from Go, keep alive period + idle are ~15sec. The default
	// keep count is 8 according to some kernel docs. That means it should
	// fail after ~120 seconds. Unfortunately, keep alive only happens if
	// there is no send-q on the socket, and so a slow reader can still cause
	// hanging sockets forever. By setting user timeout, we will kill the
	// connection if any writes go unacknowledged for the amount of time.
	// This should close the keep alive hole.
	//
	// See https://blog.cloudflare.com/when-tcp-sockets-refuse-to-die/

	rawConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	controlErr := rawConn.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_TCP, unix.TCP_USER_TIMEOUT, int(timeout.Milliseconds()))
	})
	if controlErr != nil {
		return controlErr
	}
	if err != nil {
		return err
	}
	return nil
}
