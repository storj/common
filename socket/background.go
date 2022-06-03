// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package socket

import (
	"net"
	"syscall"

	"github.com/zeebo/errs"
)

// BackgroundDialer returns a net.Dialer configured to try
// to set the lowest priority socket settings, changing the
// congestion controller to a background congestion controller if
// possible or available. On Linux, will use the kernel module
// specified by STORJ_SOCKET_LOWPRIO_CTL. On Linux, 'cdg' is
// recommended, with module parameters use_shadow=0 and use_ineff=0.
//
// This is useful for configuring storj.io/uplink.Config.DialContext
// to be background sockets. Expected usage like:
//
//   cfg := uplink.Config{
//     DialContext: socket.BackgroundDialer().DialContext,
//   }
//
func BackgroundDialer() *net.Dialer {
	return &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			var eg errs.Group
			eg.Add(c.Control(func(fd uintptr) {
				eg.Add(setLowPrioCongestionController(int(fd)))
				eg.Add(setLowEffortQoS(int(fd)))
			}))
			err := eg.Err()
			if err != nil {
				// should we log this?
				_ = err
			}
			return nil
		},
	}
}
