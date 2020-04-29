// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package fpath_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/fpath"
	"storj.io/common/testcontext"
)

func TestAtomicWriteFile(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	err := fpath.AtomicWriteFile(ctx.File("example.txt"), []byte{1, 2, 3}, 0600)
	require.NoError(t, err)

	data, err := ioutil.ReadFile(ctx.File("example.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3}, data)
}
