// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

/*
Package geoip provides utility code used across storj projects for leveraging the Maxmind GeoIP datasets. The maxminddb
package was not added as a dependency as this utility code is not needed in every project. Projects consuming the
database should include it in their dependencies. Example usage:

reader, err := maxminddb.Open("path/to/GeoIP2Lite.mmdb")
if err != nil {
	return errs.New("unable to open geolocation db: %w", err)
}
peer.IPDB = geoip.NewIPDB(reader)
*/
package geoip
