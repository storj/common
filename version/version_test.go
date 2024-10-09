// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package version_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
	"storj.io/common/testrand"
	"storj.io/common/version"
)

func TestInfo_IsZero(t *testing.T) {
	zeroInfo := version.Info{}
	require.True(t, zeroInfo.IsZero())

	ver, err := version.NewSemVer("1.2.3")
	require.NoError(t, err)

	info := version.Info{
		Version: ver,
	}
	require.False(t, info.IsZero())
}

func TestSemVer_String(t *testing.T) {
	buildVer, err := version.NewSemVer("1.2.3-rc")
	require.NoError(t, err)
	require.Equal(t, "v1.2.3-rc", buildVer.String())

	buildVer2, err := version.NewSemVer("1.2.3-rc-metainfo")
	require.NoError(t, err)
	require.Equal(t, "v1.2.3-rc-metainfo", buildVer2.String())

	nonBuildVer, err := version.NewSemVer("1.2.3")
	require.NoError(t, err)
	require.Equal(t, "v1.2.3", nonBuildVer.String())
}

func TestSemVer_IsZero(t *testing.T) {
	zeroVer := version.SemVer{}
	require.True(t, zeroVer.IsZero())

	ver, err := version.NewSemVer("1.2.3")
	require.NoError(t, err)
	require.False(t, ver.IsZero())
}

func TestSemVer_Compare(t *testing.T) {
	version001, err := version.NewSemVer("v0.0.1")
	require.NoError(t, err)
	version002, err := version.NewSemVer("v0.0.2")
	require.NoError(t, err)
	version030, err := version.NewSemVer("v0.3.0")
	require.NoError(t, err)
	version040, err := version.NewSemVer("v0.4.0")
	require.NoError(t, err)
	version500, err := version.NewSemVer("v5.0.0")
	require.NoError(t, err)
	version600, err := version.NewSemVer("v6.0.0")
	require.NoError(t, err)

	// compare the same values
	require.True(t, version001.Compare(version001) == 0) //nolint: gocritic
	require.True(t, version030.Compare(version030) == 0) //nolint: gocritic
	require.True(t, version500.Compare(version500) == 0) //nolint: gocritic

	require.True(t, version001.Compare(version002) < 0)
	require.True(t, version030.Compare(version040) < 0)
	require.True(t, version500.Compare(version600) < 0)
	require.True(t, version001.Compare(version030) < 0)
	require.True(t, version030.Compare(version500) < 0)

	require.True(t, version002.Compare(version001) > 0)
	require.True(t, version040.Compare(version030) > 0)
	require.True(t, version600.Compare(version500) > 0)
	require.True(t, version030.Compare(version002) > 0)
	require.True(t, version600.Compare(version040) > 0)
}

func TestVersion_IsZero(t *testing.T) {
	zeroVer := version.Version{}
	require.True(t, zeroVer.IsZero())

	ver := version.Version{Version: "v1.2.3", URL: "http://127.0.0.1/"}
	require.False(t, ver.IsZero())
}

func TestRollout_MarshalJSON_UnmarshalJSON(t *testing.T) {
	var arbitraryRollout version.Rollout
	for i := 0; i < len(version.RolloutBytes{}); i++ {
		arbitraryRollout.Seed[i] = byte(i)
		arbitraryRollout.Cursor[i] = byte(i * 2)
	}

	scenarios := []struct {
		name    string
		rollout version.Rollout
	}{
		{
			"arbitrary rollout",
			arbitraryRollout,
		},
		{
			"empty rollout",
			version.Rollout{},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			var actualRollout version.Rollout

			_, err := json.Marshal(actualRollout.Seed)
			require.NoError(t, err)

			jsonRollout, err := json.Marshal(scenario.rollout)
			require.NoError(t, err)

			err = json.Unmarshal(jsonRollout, &actualRollout)
			require.NoError(t, err)
			require.Equal(t, scenario.rollout, actualRollout)
		})
	}
}

func TestShouldUpdate(t *testing.T) {
	// NB: total and acceptable tolerance are negatively correlated.
	total := 20000
	tolerance := total / 100 // 1%

	for p := 10; p < 100; p += 10 {
		var rollouts int
		percentage := p
		cursor := version.PercentageToCursor(percentage)

		rollout := version.Rollout{
			Seed:   version.RolloutBytes{},
			Cursor: cursor,
		}
		testrand.Read(rollout.Seed[:])

		for i := 0; i < total; i++ {
			var nodeID storj.NodeID
			testrand.Read(nodeID[:])

			if version.ShouldUpdate(rollout, nodeID) {
				rollouts++
			}
		}

		assert.Condition(t, func() bool {
			diff := rollouts - (total * percentage / 100)
			return int(math.Abs(float64(diff))) < tolerance
		})
	}
}

func TestPercentageToCursorF(t *testing.T) {
	type args struct {
		percentage float64
		expected   string
	}

	tests := []args{
		{
			percentage: 0,
			expected:   "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			percentage: 6,
			expected:   "0f5c28f5c28f5c28f5c28f5c28f5c28f5c28f5c28f5c28f5c28f5c28f5c1f960",
		},
		{
			percentage: 12,
			expected:   "1eb851eb851eb851eb851eb851eb851eb851eb851eb851eb851eb851eb83f2c0",
		},
		{
			percentage: 25,
			expected:   "3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd8f10",
		},
		{
			percentage: 50,
			expected:   "7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb1e20",
		},
		{
			percentage: 100,
			expected:   "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("cursor-%.0f-percent", tt.percentage), func(t *testing.T) {
			res := version.PercentageToCursorF(tt.percentage)
			assert.Equal(t, tt.expected, hex.EncodeToString(res[:]))
		})
	}
}

func TestPercentageToCursorF_Precision(t *testing.T) {
	// check if we have enough precision to have small difference between 0.2% and 0
	cursorSmall := version.PercentageToCursorF(0.002)
	cursorZero := version.PercentageToCursorF(0)
	assert.NotEqual(t, cursorZero, cursorSmall)
}
