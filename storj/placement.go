// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import (
	"storj.io/common/storj/location"
)

// PlacementConstraint is the ID of the placement/geofencing rule.
type PlacementConstraint int

const (

	// EveryCountry includes all countries.
	EveryCountry PlacementConstraint = 0

	// EU includes only the 27 members of European Union.
	EU = 1

	// EEA defines the European Economic Area (EU + 3 countries), the area where GDPR is valid.
	EEA = 2

	// US filters nodes only from the United States.
	US = 3

	// DE placement uses nodes only from Germany.
	DE = 4
)

// AllowedCountry checks if country is allowed by the placement policy.
func (p PlacementConstraint) AllowedCountry(isoCountryCode location.CountryCode) bool {
	if p == EveryCountry {
		return true
	}
	switch p {
	case EEA:
		for _, c := range location.EuCountries {
			if c == isoCountryCode {
				return true
			}
		}
		for _, c := range location.EeaNonEuCountries {
			if c == isoCountryCode {
				return true
			}
		}
	case EU:
		for _, c := range location.EuCountries {
			if c == isoCountryCode {
				return true
			}
		}
	case US:
		return isoCountryCode.Equal(location.UnitedStates)
	case DE:
		return isoCountryCode.Equal(location.Germany)
	default:
		return false
	}
	return false
}
