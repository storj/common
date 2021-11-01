// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package testrand_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/testrand"
)

var stringLength = 50

func TestNoPanic(t *testing.T) {
	t.Log("URLPath", testrand.URLPath())
	t.Log("URLPathNonFolder", testrand.URLPathNonFolder())
	t.Log("Path", testrand.Path())
	isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString
	randomNumericString := string(testrand.RandNumeric(stringLength))
	t.Log("Numeric String", randomNumericString)
	require.True(t, isNumeric(randomNumericString))
	require.Equal(t, stringLength, len(randomNumericString))
	isAlphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	randomAlphaNumericString := string(testrand.RandAlphaNumeric(stringLength))
	t.Log("AlphaNumeric String", randomAlphaNumericString)
	require.True(t, isAlphaNumeric(randomAlphaNumericString))
	require.Equal(t, stringLength, len(randomAlphaNumericString))
}
