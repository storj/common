// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import (
	"database/sql/driver"

	"github.com/zeebo/errs"

	"storj.io/common/storj/location"
)

// PlacementConstraint is the ID of the placement/geofencing rule.
type PlacementConstraint uint16

const (
	// EveryCountry includes all countries.
	EveryCountry PlacementConstraint = 0

	// EU includes only the 27 members of European Union.
	EU PlacementConstraint = 1

	// EEA defines the European Economic Area (EU + 3 countries), the area where GDPR is valid.
	EEA PlacementConstraint = 2

	// US filters nodes only from the United States.
	US PlacementConstraint = 3

	// DE placement uses nodes only from Germany.
	DE PlacementConstraint = 4

	// InvalidPlacement is used when there is no information about the stored placement.
	InvalidPlacement PlacementConstraint = 5

	// NR placement uses nodes that are not in RU or other countries sanctioned because of the RU/UA War.
	NR PlacementConstraint = 6
)

var placementConstraintLookup = [...]location.Set{
	EveryCountry: location.NewFullSet(),
	EU:           location.EuCountries,
	EEA:          location.EeaCountries,
	US:           location.NewSet(location.UnitedStates),
	DE:           location.NewSet(location.Germany),
	NR:           location.NewFullSet().Without(location.Russia, location.Belarus),
}

// AllowedCountry checks if country is allowed by the placement policy.
func (p PlacementConstraint) AllowedCountry(isoCountryCode location.CountryCode) bool {
	if int(p) >= len(placementConstraintLookup) {
		return false
	}
	return placementConstraintLookup[p].Contains(isoCountryCode)
}

// Value implements the driver.Valuer interface.
func (p PlacementConstraint) Value() (driver.Value, error) {
	return int64(p), nil
}

// Scan implements the sql.Scanner interface.
func (p *PlacementConstraint) Scan(value interface{}) error {
	if value == nil {
		*p = EveryCountry
		return nil
	}

	if _, isInt64 := value.(int64); !isInt64 {
		return errs.New("unable to scan %T into PlacementConstraint", value)
	}

	code, err := driver.Int32.ConvertValue(value)
	if err != nil {
		return errs.Wrap(err)
	}
	*p = PlacementConstraint(uint16(code.(int64)))
	return nil

}
