// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package socket

func setLowPrioCongestionController(fd int) error { return nil }

func setLowEffortQoS(fd int) error { return nil }
