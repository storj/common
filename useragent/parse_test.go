// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package useragent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/useragent"
)

func TestParseEntries(t *testing.T) {
	type test struct {
		in  string
		exp []useragent.Entry
	}

	var tests = []test{{
		in:  ``,
		exp: []useragent.Entry{},
	}, {
		in: `  Mozilla`,
		exp: []useragent.Entry{
			{"Mozilla", "", ""},
		},
	}, {
		in: `Mozilla   `,
		exp: []useragent.Entry{
			{"Mozilla", "", ""},
		},
	}, {
		in: `Mozilla`,
		exp: []useragent.Entry{
			{"Mozilla", "", ""},
		},
	}, {
		in: `Mozilla/5.0`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
		},
	}, {
		in: `Mozilla/5.0 (Linux; U; Android 4.4.3;)`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
		},
	}, {
		in: `Mozilla/5.0 (Linux; \(U\); Android 4.4.3;)`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; (U); Android 4.4.3;"},
		},
	}, {
		in: `Mozilla/5.0 (Linux; U; Android 4.4.3;) Mobile`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
			{"Mobile", "", ""},
		},
	}, {
		in: `Mozilla/5.0 (Linux; U; Android 4.4.3;) Mobile Safari/534.30`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
			{"Mobile", "", ""},
			{"Safari", "534.30", ""},
		},
	}, {
		in: `storj.io-uplink/v0.0.1`,
		exp: []useragent.Entry{
			{"storj.io-uplink", "v0.0.1", ""},
		},
	}, {
		in: `storj.io-uplink/v0.0.1 storj.io-drpc/v5.0.0+123+123`,
		exp: []useragent.Entry{
			{"storj.io-uplink", "v0.0.1", ""},
			{"storj.io-drpc", "v5.0.0+123+123", ""},
		},
	}, {
		in: `Mozilla/5.0 (Linux; U; Android 4.4.3;) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 Opera News/1.0`,
		exp: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
			{"AppleWebkit", "534.30", ""},
			{"", "", "KHTML, like Gecko"},
			{"Version", "4.0", ""},
			{"Mobile", "", ""},
			{"Safari", "534.30", ""},
			{"Opera", "", ""},
			{"News", "1.0", ""},
		},
	}}

	for _, test := range tests {
		entries, err := useragent.ParseEntries([]byte(test.in))
		if !assert.NoError(t, err, test.in) {
			continue
		}
		assert.Equal(t, test.exp, entries, test.in)
	}
}

func TestParseInvalid(t *testing.T) {
	type test struct {
		in string
	}

	var tests = []test{
		// comment must not be first
		{`(Linux) Mozilla`},
		// invalid comments
		{`Mozilla (Li ( nux)`},
		{`Mozilla (Li ) nux)`},
		// although valid per RFC, it's unsupported for now
		{`Mozilla/5.0 (Linux; (U; Android) 4.4.3;)`},
	}

	for _, test := range tests {
		_, err := useragent.ParseEntries([]byte(test.in))
		assert.Error(t, err, test.in)
	}
}
