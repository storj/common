// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testcontext_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/testcontext"
)

func TestCompile(t *testing.T) {
	t.Skip("temporarily disabled")
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	exe := ctx.Compile("storj.io/storj/examples/grpc-debug")
	assert.NotEmpty(t, exe)
}
