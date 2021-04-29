// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package httpranger implements lazy io.Reader and io.Writer interfaces.
package httpranger

import (
	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
)

// Error is the errs class of standard Ranger errors.
var Error = errs.Class("ranger")

var mon = monkit.Package()
