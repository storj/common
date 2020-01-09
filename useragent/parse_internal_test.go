// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package useragent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	type test struct {
		in    string
		token string
		next  int
		ok    bool
	}

	tests := []test{
		{``, ``, 0, false},
		{` `, ``, 0, false},
		{`(`, ``, 0, false},
		{`)`, ``, 0, false},
		{`a`, `a`, 1, true},
		{`a `, `a`, 1, true},
		{`a b`, `a`, 1, true},
		{`a/x b`, `a`, 1, true},
	}

	for _, test := range tests {
		token, next, ok := parseToken([]byte(test.in), 0)
		assert.Equal(t, test.token, token, test.in)
		assert.Equal(t, test.next, next, test.in)
		assert.Equal(t, test.ok, ok, test.in)
	}
}
