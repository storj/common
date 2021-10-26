// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package storj

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/storj/location"
)

func TestCriteria_Geofencing(t *testing.T) {

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
			assert.Equal(t, c.expected, c.placement.allowedCountry(c.country))
		})
	}
}
