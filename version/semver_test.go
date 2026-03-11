// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		input   string
		major   uint64
		minor   uint64
		patch   uint64
		pre     string
		build   string
		str     string
		wantErr bool
	}{
		// Basic versions.
		{input: "1.2.3", major: 1, minor: 2, patch: 3, str: "1.2.3"},
		{input: "0.0.0", major: 0, minor: 0, patch: 0, str: "0.0.0"},
		{input: "10.20.30", major: 10, minor: 20, patch: 30, str: "10.20.30"},

		// With v prefix.
		{input: "v1.2.3", major: 1, minor: 2, patch: 3, str: "1.2.3"},
		{input: "v0.0.1", major: 0, minor: 0, patch: 1, str: "0.0.1"},

		// Storj classic: v1.119.3-commithash.
		{input: "v1.119.3-abcdef1234", major: 1, minor: 119, patch: 3, pre: "abcdef1234", str: "1.119.3-abcdef1234"},
		{input: "v1.123.2-5972e87e2", major: 1, minor: 123, patch: 2, pre: "5972e87e2", str: "1.123.2-5972e87e2"},

		// Storj classic with commit hash starting with zero (the bug case).
		{input: "v1.119.3-093655760", major: 1, minor: 119, patch: 3, pre: "093655760", str: "1.119.3-093655760"},
		{input: "v1.119.3-00001234", major: 1, minor: 119, patch: 3, pre: "00001234", str: "1.119.3-00001234"},

		// Storj date-based: v2026.01.1769691272-5972e87e2.
		{input: "v2026.01.1769691272-5972e87e2", major: 2026, minor: 1, patch: 1769691272, pre: "5972e87e2", str: "2026.1.1769691272-5972e87e2"},
		{input: "v2026.03.1234567890-0abcdef12", major: 2026, minor: 3, patch: 1234567890, pre: "0abcdef12", str: "2026.3.1234567890-0abcdef12"},

		// Pre-release with multiple dot-separated identifiers.
		{input: "1.2.3-rc.1", major: 1, minor: 2, patch: 3, pre: "rc.1", str: "1.2.3-rc.1"},
		{input: "1.2.3-alpha.0.beta", major: 1, minor: 2, patch: 3, pre: "alpha.0.beta", str: "1.2.3-alpha.0.beta"},

		// Pre-release with hyphen-separated identifiers (single identifier with hyphens).
		{input: "1.2.3-rc-metainfo", major: 1, minor: 2, patch: 3, pre: "rc-metainfo", str: "1.2.3-rc-metainfo"},

		// Build metadata.
		{input: "1.2.3+build123", major: 1, minor: 2, patch: 3, build: "build123", str: "1.2.3+build123"},
		{input: "1.2.3-rc.1+build", major: 1, minor: 2, patch: 3, pre: "rc.1", build: "build", str: "1.2.3-rc.1+build"},

		// Tolerant parsing: leading zeroes in major/minor/patch stripped.
		{input: "01.02.03", major: 1, minor: 2, patch: 3, str: "1.2.3"},

		// Tolerant parsing: shortened versions filled.
		{input: "1.2", major: 1, minor: 2, patch: 0, str: "1.2.0"},
		{input: "1", major: 1, minor: 0, patch: 0, str: "1.0.0"},

		// Whitespace trimming.
		{input: "  v1.2.3  ", major: 1, minor: 2, patch: 3, str: "1.2.3"},

		// Errors.
		{input: "", wantErr: true},
		{input: "not-a-version", wantErr: true},
		{input: "1.2.3-", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, err := ParseSemVer(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.major, v.Major)
			assert.Equal(t, tt.minor, v.Minor)
			assert.Equal(t, tt.patch, v.Patch)
			assert.Equal(t, tt.str, v.String())

			// Verify pre-release round-trips correctly.
			if tt.pre != "" {
				var preParts []string
				for _, p := range v.Pre {
					preParts = append(preParts, p.String())
				}
				got := ""
				for i, p := range preParts {
					if i > 0 {
						got += "."
					}
					got += p
				}
				assert.Equal(t, tt.pre, got)
			}

			// Verify build round-trips correctly.
			if tt.build != "" {
				got := ""
				for i, b := range v.Build {
					if i > 0 {
						got += "."
					}
					got += b
				}
				assert.Equal(t, tt.build, got)
			}
		})
	}
}

func TestParseSemVer_LeadingZeroPrerelease(t *testing.T) {
	// This is the exact case that caused the panic with blang/semver.
	v, err := ParseSemVer("v1.119.3-093655760")
	require.NoError(t, err)
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(119), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	require.Len(t, v.Pre, 1)
	// Should be preserved as-is (alphanumeric, not numeric).
	assert.False(t, v.Pre[0].IsNumeric())
	assert.Equal(t, "093655760", v.Pre[0].String())
}

func TestSemverCompare(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		// Basic ordering.
		{"1.0.0", "2.0.0", -1},
		{"1.0.0", "1.1.0", -1},
		{"1.0.0", "1.0.1", -1},
		{"1.0.0", "1.0.0", 0},
		{"2.0.0", "1.0.0", 1},

		// Pre-release has lower precedence than release.
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha", 1},

		// Pre-release ordering.
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-1", "1.0.0-2", -1},
		{"1.0.0-1", "1.0.0-alpha", -1}, // numeric < alphanumeric

		// Storj versions: same base, different commit hash.
		{"v1.119.3-aaa", "v1.119.3-bbb", -1},
		{"v1.119.3-bbb", "v1.119.3-aaa", 1},

		// Storj versions: different patch.
		{"v1.119.2-abc", "v1.119.3-abc", -1},

		// Date-based versions.
		{"v2026.01.100-abc", "v2026.01.200-abc", -1},
		{"v2026.01.100-abc", "v2026.02.100-abc", -1},

		// Leading zero prerelease is alphanumeric, non-leading-zero is numeric.
		// Alphanumeric > numeric in semver precedence.
		{"v1.0.0-093655760", "v1.0.0-193655760", 1},
		// Two leading-zero prereleases compare lexically.
		{"v1.0.0-01234", "v1.0.0-05678", -1},

		// More pre-release identifiers means greater if all preceding match.
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			a, err := ParseSemVer(tt.a)
			require.NoError(t, err)
			b, err := ParseSemVer(tt.b)
			require.NoError(t, err)
			assert.Equal(t, tt.want, a.Compare(b))
		})
	}
}

func TestSemverCompareHelpers(t *testing.T) {
	a, err := ParseSemVer("1.0.0")
	require.NoError(t, err)
	b, err := ParseSemVer("2.0.0")
	require.NoError(t, err)

	assert.True(t, a.LT(b))
	assert.True(t, a.LE(b))
	assert.False(t, a.GT(b))
	assert.False(t, a.GE(b))
	assert.False(t, a.EQ(b))

	assert.True(t, a.EQ(a))
	assert.True(t, a.LE(a))
	assert.True(t, a.GE(a))
}

func TestSemverJSON(t *testing.T) {
	type wrapper struct {
		Version SemVer `json:"version"`
	}

	original := wrapper{}
	var err error
	original.Version, err = ParseSemVer("v1.119.3-093655760")
	require.NoError(t, err)

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.Contains(t, string(data), "1.119.3-093655760")

	var decoded wrapper
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, original.Version.String(), decoded.Version.String())
	assert.Equal(t, 0, original.Version.Compare(decoded.Version))
}

func TestSemverJSON_DateBased(t *testing.T) {
	type wrapper struct {
		Version SemVer `json:"version"`
	}

	original := wrapper{}
	var err error
	original.Version, err = ParseSemVer("v2026.01.1769691272-5972e87e2")
	require.NoError(t, err)

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded wrapper
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Equal(t, original.Version.String(), decoded.Version.String())
}

func TestPrereleaseVersionCompare(t *testing.T) {
	num1, err := NewPrereleaseVersion("1")
	require.NoError(t, err)
	num2, err := NewPrereleaseVersion("2")
	require.NoError(t, err)
	alpha, err := NewPrereleaseVersion("alpha")
	require.NoError(t, err)
	beta, err := NewPrereleaseVersion("beta")
	require.NoError(t, err)
	leading, err := NewPrereleaseVersion("093655760")
	require.NoError(t, err)

	// Numeric vs numeric.
	assert.Equal(t, -1, num1.Compare(num2))
	assert.Equal(t, 1, num2.Compare(num1))
	num1b, err := NewPrereleaseVersion("1")
	require.NoError(t, err)
	assert.Equal(t, 0, num1.Compare(num1b))

	// Alpha vs alpha.
	assert.Equal(t, -1, alpha.Compare(beta))
	assert.Equal(t, 1, beta.Compare(alpha))

	// Numeric < alpha.
	assert.Equal(t, -1, num1.Compare(alpha))
	assert.Equal(t, 1, alpha.Compare(num1))

	// Leading zero is alphanumeric.
	assert.False(t, leading.IsNumeric())
	assert.Equal(t, "093655760", leading.String())
}

func FuzzParseSemVer(f *testing.F) {
	f.Add("1.2.3")
	f.Add("v1.119.3-093655760")
	f.Add("v2026.01.1769691272-5972e87e2")
	f.Add("1.2.3-rc.1+build123")
	f.Add("0.0.0")
	f.Add("v1.2")
	f.Add("01.02.03")
	f.Add("")
	f.Add("not-a-version")

	f.Fuzz(func(t *testing.T, input string) {
		v, err := ParseSemVer(input)
		if err != nil {
			return
		}
		// Round-trip: parse the String() output and compare.
		v2, err := ParseSemVer(v.String())
		if err != nil {
			t.Fatalf("failed to re-parse %q (from %q): %v", v.String(), input, err)
		}
		if v.Compare(v2) != 0 {
			t.Fatalf("round-trip mismatch: %q -> %q -> %q", input, v.String(), v2.String())
		}
	})
}
