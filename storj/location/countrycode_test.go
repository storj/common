// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package location

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountryCode_String(t *testing.T) {
	require.Equal(t, "HU", ToCountryCode("HU").String())
	require.Equal(t, "DE", ToCountryCode("DE").String())
	require.Equal(t, "XX", ToCountryCode("XX").String())
	require.Equal(t, "", None.String())
}
