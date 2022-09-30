// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package hmacsha512

import "testing"

func TestGenericBlock(t *testing.T) {
	// This is here to avoid linters complaining about blockGeneric being unused.
	var d digest
	d.Reset()
	var block [BlockSize]byte
	blockGeneric(&d, block[:])
}
