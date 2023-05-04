// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package leak

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResource_Ok(t *testing.T) {
	root := RootResource(0)
	child1 := root.Child("alpha", 0)
	child2 := root.Child("beta", 0)
	require.NoError(t, child1.Close())
	require.NoError(t, child2.Close())
	require.NoError(t, root.Close())
}

func TestResource_Nested(t *testing.T) {
	root := RootResource(0)
	child1 := root.Child("alpha", 0)
	leak2 := root.Child("beta", 0)
	require.NoError(t, child1.Close())
	_ = leak2
	require.Error(t, root.Close())
}
