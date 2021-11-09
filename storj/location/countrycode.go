// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package location

import "strings"

// CountryCode stores upper case ISO code of countries.
type CountryCode uint16

// ToCountryCode convert string to CountryCode.
// encoding is based on the ASCII representation of the country code.
func ToCountryCode(s string) CountryCode {
	if len(s) != 2 {
		return CountryCode(0)
	}
	upper := strings.ToUpper(s)
	return CountryCode(uint16(upper[0])*uint16(256) + uint16(upper[1]))
}

// Equal compares two country code.
func (c CountryCode) Equal(o CountryCode) bool {
	return c == o
}

// String returns with the upper-case (two letter) ISO code of the country.
func (c CountryCode) String() string {
	if c == 0 {
		return ""
	}
	return string([]byte{byte(c / 256), byte(c % 256)})
}
