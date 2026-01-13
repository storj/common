// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"storj.io/common/version/buildinfo"
)

// FromBuild returns version string for a module.
//
// Deprecated: use buildinfo package, which doesn't have any 3rd party dependencies.
func FromBuild(modname string) (string, error) {
	return buildinfo.FromBuild(modname)
}
