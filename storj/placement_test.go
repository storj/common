// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package storj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlacement_SQLConversion(t *testing.T) {
	p := PlacementConstraint(2)
	value, err := p.Value()
	require.NoError(t, err)

	res := new(PlacementConstraint)
	err = res.Scan(value)
	require.NoError(t, err)
	require.Equal(t, PlacementConstraint(2), *res)

	err = res.Scan(nil)
	require.NoError(t, err)
	require.Equal(t, DefaultPlacement, *res)

	err = res.Scan("")
	require.Error(t, err)
}
