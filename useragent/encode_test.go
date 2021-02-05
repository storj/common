// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package useragent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/useragent"
)

func TestEncodeEntries(t *testing.T) {
	// invalid product
	_, err := useragent.EncodeEntries([]useragent.Entry{
		{"Prod)uct", "", ""},
	})
	require.Error(t, err)

	// invalid version
	_, err = useragent.EncodeEntries([]useragent.Entry{
		{"Product", "Vers(ion", ""},
	})
	require.Error(t, err)

	type test struct {
		in  []useragent.Entry
		exp string
	}

	var tests = []test{{
		in: []useragent.Entry{
			{"Mozilla", "", ""},
		},
		exp: `Mozilla`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", ""},
		},
		exp: `Mozilla/5.0`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", "Linux; U; Android 4.4.3;"},
		},
		exp: `Mozilla/5.0 (Linux; U; Android 4.4.3;)`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "", "Linux; U; Android 4.4.3;"},
		},
		exp: `Mozilla (Linux; U; Android 4.4.3;)`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
		},
		exp: `Mozilla/5.0 (Linux; U; Android 4.4.3;)`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; (U); Android 4.4.3;"},
		},
		exp: `Mozilla/5.0 (Linux; (U); Android 4.4.3;)`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
			{"Mobile", "", ""},
		},
		exp: `Mozilla/5.0 (Linux; U; Android 4.4.3;) Mobile`,
	}, {
		in: []useragent.Entry{
			{"Mozilla", "5.0", ""},
			{"", "", "Linux; U; Android 4.4.3;"},
			{"Mobile", "", ""},
			{"Safari", "534.30", ""},
		},
		exp: `Mozilla/5.0 (Linux; U; Android 4.4.3;) Mobile Safari/534.30`,
	}, {
		in: []useragent.Entry{
			{"storj.io-uplink", "v0.0.1", ""},
		},
		exp: `storj.io-uplink/v0.0.1`,
	}, {
		in: []useragent.Entry{
			{"storj.io-uplink", "v0.0.1", ""},
			{"storj.io-drpc", "v5.0.0+123+123", ""},
		},
		exp: `storj.io-uplink/v0.0.1 storj.io-drpc/v5.0.0+123+123`,
	}, {
		in: []useragent.Entry{
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
		exp: `Mozilla/5.0 (Linux; U; Android 4.4.3;) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 Opera News/1.0`,
	}, {
		in: []useragent.Entry{
			{"Blocknify", "", ""},
			{"Uplink", "1.4.6-0.20210201122710-48b82ce14a37", ""},
		},
		exp: `Blocknify Uplink/1.4.6-0.20210201122710-48b82ce14a37`,
	}}

	for _, test := range tests {
		encoded, err := useragent.EncodeEntries(test.in)
		require.NoError(t, err)
		assert.Equal(t, test.exp, string(encoded))
	}
}
