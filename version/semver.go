// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

// SemVer represents a parsed semantic version.
type SemVer struct {
	Major uint64
	Minor uint64
	Patch uint64
	Pre   []PrereleaseVersion
	Build []string
}

// PrereleaseVersion represents a single pre-release identifier,
// which is either numeric or alphanumeric.
type PrereleaseVersion struct {
	versionStr string
	versionNum uint64
	isNum      bool
}

// NewPrereleaseVersion creates a pre-release version from a string.
func NewPrereleaseVersion(s string) (PrereleaseVersion, error) {
	if s == "" {
		return PrereleaseVersion{}, fmt.Errorf("empty prerelease version")
	}
	// If it's all digits and has no leading zero, treat as numeric.
	if isAllDigits(s) && !hasLeadingZero(s) {
		num, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return PrereleaseVersion{}, err
		}
		return PrereleaseVersion{versionNum: num, isNum: true}, nil
	}
	return PrereleaseVersion{versionStr: s, isNum: false}, nil
}

// IsNumeric returns true if the pre-release identifier is numeric.
func (v PrereleaseVersion) IsNumeric() bool {
	return v.isNum
}

// Compare compares two pre-release versions.
// Per semver spec: numeric identifiers are compared as integers,
// alphanumeric identifiers are compared lexically,
// numeric identifiers always have lower precedence than alphanumeric.
func (v PrereleaseVersion) Compare(o PrereleaseVersion) int {
	if v.isNum && o.isNum {
		return cmp.Compare(v.versionNum, o.versionNum)
	}
	if v.isNum {
		return -1 // numeric < alphanumeric
	}
	if o.isNum {
		return 1 // alphanumeric > numeric
	}
	return strings.Compare(v.versionStr, o.versionStr)
}

// String returns the string representation of the pre-release version.
func (v PrereleaseVersion) String() string {
	if v.isNum {
		return strconv.FormatUint(v.versionNum, 10)
	}
	return v.versionStr
}

// Compare compares two versions.
func (v SemVer) Compare(o SemVer) int {
	return cmp.Or(
		cmp.Compare(v.Major, o.Major),
		cmp.Compare(v.Minor, o.Minor),
		cmp.Compare(v.Patch, o.Patch),
		comparePre(v.Pre, o.Pre),
	)
}

// Equal returns true if v equals o.
func (v SemVer) Equal(o SemVer) bool {
	return v.Compare(o) == 0
}

// EQ returns true if v equals o.
//
//go:fix inline
//nolint:gocheckcompilerdirectives
func (v SemVer) EQ(o SemVer) bool {
	return v.Equal(o)
}

// Less returns true if v is less than o.
func (v SemVer) Less(o SemVer) bool {
	return v.Compare(o) < 0
}

// LT returns true if v is less than o.
//
//go:fix inline
//nolint:gocheckcompilerdirectives
func (v SemVer) LT(o SemVer) bool {
	return v.Less(o)
}

// Greater returns true if v is greater than o.
func (v SemVer) Greater(o SemVer) bool {
	return v.Compare(o) > 0
}

// GT returns true if v is greater than o.
//
//go:fix inline
//nolint:gocheckcompilerdirectives
func (v SemVer) GT(o SemVer) bool {
	return v.Greater(o)
}

// LessOrEqual returns true if v is less than or equal to o.
func (v SemVer) LessOrEqual(o SemVer) bool {
	return v.Compare(o) <= 0
}

// LE returns true if v is less than or equal to o.
//
//go:fix inline
//nolint:gocheckcompilerdirectives
func (v SemVer) LE(o SemVer) bool {
	return v.LessOrEqual(o)
}

// GreaterOrEqual returns true if v is greater than or equal to o.
func (v SemVer) GreaterOrEqual(o SemVer) bool {
	return v.Compare(o) >= 0
}

// GE returns true if v is greater than or equal to o.
//
//go:fix inline
//nolint:gocheckcompilerdirectives
func (v SemVer) GE(o SemVer) bool {
	return v.GreaterOrEqual(o)
}

// VString returns the semver string with a "v" prefix.
func (v SemVer) VString() string {
	return "v" + v.String()
}

// String returns the semver string without the "v" prefix.
func (v SemVer) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.Pre) > 0 {
		b.WriteByte('-')
		for i, p := range v.Pre {
			if i > 0 {
				b.WriteByte('.')
			}
			b.WriteString(p.String())
		}
	}
	if len(v.Build) > 0 {
		b.WriteByte('+')
		for i, s := range v.Build {
			if i > 0 {
				b.WriteByte('.')
			}
			b.WriteString(s)
		}
	}
	return b.String()
}

// IsZero checks if the semantic version is its zero value.
func (v SemVer) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && len(v.Pre) == 0 && len(v.Build) == 0
}

// MarshalText implements encoding.TextMarshaler.
func (v SemVer) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (v *SemVer) UnmarshalText(data []byte) error {
	parsed, err := ParseSemVer(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// ParseSemVer parses a version string tolerantly.
func ParseSemVer(s string) (SemVer, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")

	if s == "" {
		return SemVer{}, fmt.Errorf("empty version string")
	}

	var v SemVer

	// Split off build metadata first (after +).
	if before, after, ok := strings.Cut(s, "+"); ok {
		v.Build = strings.Split(after, ".")
		s = before
	}

	// Split off pre-release (after first -).
	var preStr string
	if before, after, ok := strings.Cut(s, "-"); ok {
		preStr = after
		s = before
		if preStr == "" {
			return SemVer{}, fmt.Errorf("empty prerelease version")
		}
	}

	// Parse major.minor.patch (tolerant: allows 1 or 2 components).
	parts := strings.SplitN(s, ".", 3)

	// Validate that parts look like numbers.
	for _, p := range parts {
		if p == "" || !isAllDigits(p) {
			return SemVer{}, fmt.Errorf("invalid version component: %q", p)
		}
	}

	var err error
	v.Major, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid major version: %w", err)
	}

	if len(parts) >= 2 {
		v.Minor, err = strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			return SemVer{}, fmt.Errorf("invalid minor version: %w", err)
		}
	}

	if len(parts) >= 3 {
		v.Patch, err = strconv.ParseUint(parts[2], 10, 64)
		if err != nil {
			return SemVer{}, fmt.Errorf("invalid patch version: %w", err)
		}
	}

	// Parse pre-release identifiers.
	if preStr != "" {
		preParts := strings.Split(preStr, ".")
		v.Pre = make([]PrereleaseVersion, len(preParts))
		for i, p := range preParts {
			v.Pre[i], err = NewPrereleaseVersion(p)
			if err != nil {
				return SemVer{}, fmt.Errorf("invalid prerelease version: %w", err)
			}
		}
	}

	return v, nil
}

func comparePre(a, b []PrereleaseVersion) int {
	aLen := len(a)
	bLen := len(b)

	if aLen == 0 && bLen == 0 {
		return 0
	}
	// Having a pre-release means lower precedence than no pre-release.
	if aLen == 0 {
		return 1
	}
	if bLen == 0 {
		return -1
	}

	minLen := aLen
	if bLen < minLen {
		minLen = bLen
	}
	for i := 0; i < minLen; i++ {
		if c := a[i].Compare(b[i]); c != 0 {
			return c
		}
	}

	// More identifiers means greater precedence if all preceding match.
	return cmp.Compare(aLen, bLen)
}

func isAllDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func hasLeadingZero(s string) bool {
	return len(s) > 1 && s[0] == '0'
}
