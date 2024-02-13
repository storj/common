// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"os"
	"storj.io/common/version/buildinfo"
)

func main() {
	versionstr, err := buildinfo.FromBuild("storj.io/common")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("%#v", versionstr)
}
