// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package storj

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/storj/location"
)

func TestPlacement_Geofencing(t *testing.T) {

	cases := []struct {
		name      string
		country   location.CountryCode
		placement PlacementConstraint
		expected  bool
	}{
		{
			name:      "US matches US selector",
			country:   location.UnitedStates,
			placement: US,
			expected:  true,
		},
		{
			name:      "Germany is EU",
			country:   location.Germany,
			placement: EU,
			expected:  true,
		},
		{
			name:      "US is not eu",
			country:   location.UnitedStates,
			placement: EU,
			expected:  false,
		},
		{
			name:      "Lower case country code is handled",
			country:   location.Germany,
			placement: EU,
			expected:  true,
		},
		{
			name:      "Empty country doesn't match region",
			country:   location.CountryCode(0),
			placement: EU,
			expected:  false,
		},
		{
			name:      "Empty country doesn't match country",
			country:   location.CountryCode(0),
			placement: US,
			expected:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.placement.AllowedCountry(c.country))
		})
	}
}

func TestPlacement_SQLConversion(t *testing.T) {
	p := PlacementConstraint(EEA)
	value, err := p.Value()
	require.NoError(t, err)

	res := new(PlacementConstraint)
	err = res.Scan(value)
	require.NoError(t, err)
	require.Equal(t, PlacementConstraint(EEA), *res)

	err = res.Scan(nil)
	require.NoError(t, err)
	require.Equal(t, EveryCountry, *res)

	err = res.Scan("")
	require.Error(t, err)
}
