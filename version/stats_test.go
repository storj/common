// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"testing"
)

func TestOsVersion(t *testing.T) {
	t.Log(osversion())
}
