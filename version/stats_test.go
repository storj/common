// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"testing"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/stretchr/testify/require"
)

func TestOsVersion(t *testing.T) {
	t.Log(osversion())
}

func TestCommitHashCRC(t *testing.T) {
	info := Info{
		CommitHash: "0e7695f391df6c2ea03d2dec02b9059ccaffb9c9",
	}

	getCommitHash := func() (ret float64) {
		info.Stats(func(key monkit.SeriesKey, field string, val float64) {
			if key.Measurement == "version_info" && field == "commit" {
				ret = val
			}
		})
		return ret
	}

	// first time the value is calculated, second time cached value is used.
	// Let's test both.
	require.Equal(t, float64(868966439), getCommitHash())
	require.Equal(t, float64(868966439), getCommitHash())
}
