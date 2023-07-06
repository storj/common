// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package location

// EuCountries defines the 27 member country of European Union.
var EuCountries = NewSet(
	Austria,
	Belgium,
	Bulgaria,
	Croatia,
	Cyprus,
	Czechia,
	Denmark,
	Estonia,
	Finland,
	France,
	Germany,
	Greece,
	Hungary,
	Ireland,
	Italy,
	Lithuania,
	Latvia,
	Luxembourg,
	Malta,
	Netherlands,
	Poland,
	Portugal,
	Romania,
	Slovenia,
	Slovakia,
	Spain,
	Sweden,
)

// EeaCountries defined the EEA countries.
var EeaCountries = func() Set {
	r := EuCountries
	r.Include(Iceland)
	r.Include(Liechtenstein)
	r.Include(Norway)
	return r
}()
