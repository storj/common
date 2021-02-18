// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testcontext_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/testcontext"
)

func TestCompile(t *testing.T) {
	ctx := testcontext.New(t)

	exe := ctx.Compile("storj.io/common/testcontext/testdata/hello")
	assert.NotEmpty(t, exe)

	exemod := ctx.CompileAt("./testdata/hellomod", "test/hello")
	assert.NotEmpty(t, exemod)
}
