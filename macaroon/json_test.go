// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package macaroon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	bs, err := json.Marshal(Caveat_Path{
		Bucket:              []byte("bucket"),
		EncryptedPathPrefix: []byte("prefix"),
	})
	require.NoError(t, err)
	require.Equal(t, string(bs), `{"bucket":"YnVja2V0","encrypted_path_prefix":"cHJlZml4"}`)
}
