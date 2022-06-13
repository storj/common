// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.16
// +build go1.16

package quic

import (
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
)

var (
	mon = monkit.Package()

	// Error is a pkg/quic error.
	Error = errs.Class("quic")

	// ErrQuicDisabled indicates QUIC has been disabled at build time.
	ErrQuicDisabled = Error.New("disabled at build time")
)
