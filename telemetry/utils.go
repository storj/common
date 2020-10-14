// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package telemetry

import (
	"math/rand"
	"net"
	"time"
)

// UnknownInstanceID is returned when no instance ID can be returned.
const UnknownInstanceID = "unknown"

// DefaultInstanceID will return the first non-nil mac address if possible, unknown otherwise.
func DefaultInstanceID() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return UnknownInstanceID
	}
	for _, iface := range ifaces {
		if iface.HardwareAddr != nil {
			return iface.HardwareAddr.String()
		}
	}
	return UnknownInstanceID
}

func jitter(t time.Duration) time.Duration {
	nanos := rand.NormFloat64()*float64(t/4) + float64(t)
	if nanos <= 0 {
		nanos = 1
	}
	return time.Duration(nanos)
}
