// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.18
// +build go1.18

package useragent_test

import (
	"testing"

	"storj.io/common/useragent"
)

func FuzzParseEntries(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("storj.io-common/v0.0.0-00010101000000-000000000000"))
	f.Add([]byte("storj.io-common/v0.0.0-00010101000000"))
	f.Add([]byte("storj.io-common/v9.0.0"))
	f.Add([]byte("Mozilla"))
	f.Add([]byte("Mozilla/5.0"))
	f.Add([]byte("Mozilla/5.0 (Linux; U; Android 4.4.3;)"))
	f.Add([]byte("storj.io-uplink/v0.0.1 storj.io-drpc/v5.0.0+123+123 Mozilla/5.0 (Linux; U; Android 4.4.3;) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30 Opera News/1.0"))
	f.Add([]byte("!#$%&'*+-.^_`|~/!#$%&'*+-.^_`|~"))

	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = useragent.ParseEntries(data)
	})
}
