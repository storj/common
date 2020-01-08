// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package testrand_test

import (
	"testing"

	"storj.io/common/testrand"
)

func TestNoPanic(t *testing.T) {
	t.Log("URLPath", testrand.URLPath())
	t.Log("URLPathNonFolder", testrand.URLPathNonFolder())
	t.Log("Path", testrand.Path())
}
