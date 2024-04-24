// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build race

package leak

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRef_Ok(t *testing.T) {
	root := Root(0)
	child1 := root.Child("alpha", 0)
	child2 := root.Child("beta", 0)
	require.NoError(t, child1.Close())
	require.NoError(t, child2.Close())
	require.NoError(t, root.Close())
}

func TestRef_Nested(t *testing.T) {
	root := Root(0)
	child1 := root.Child("alpha", 0)
	leak2 := root.Child("beta", 0)
	require.NoError(t, child1.Close())
	_ = leak2
	require.Error(t, root.Close())
}

func TestRef_Context(t *testing.T) {
	bg := context.Background()

	ref, ctx := WithContext(bg)

	root := FromContext(ctx)
	require.Equal(t, root, ref)

	child := root.Child("alpha", 0)
	require.NoError(t, child.Close())

	require.NoError(t, root.Close())
}
