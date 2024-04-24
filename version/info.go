// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.18

package version

import (
	"runtime/debug"
	"strconv"
	"time"
)

func init() {
	i := getInfoFromBuildTags()
	undefinedTime := time.Time{}
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if i.CommitHash == "" {
					i.CommitHash = s.Value
				}
			case "vcs.time":
				if i.Timestamp == undefinedTime {
					i.Timestamp, _ = time.Parse(time.RFC3339Nano, s.Value)
				}
			case "vcs.modified":
				modified, err := strconv.ParseBool(s.Value)
				if err == nil {
					i.Modified = i.Modified || modified
				}
			}
		}

	}
	Build = i
}
