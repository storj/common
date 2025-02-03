// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information.

package usermeta

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserMetaMarshal(t *testing.T) {
	meta := UserMeta{
		"key1": "value1",
		"key2": "value2",
	}

	data, err := Marshal(meta)
	require.NoError(t, err)

	meta2, err := Unmarshal(data)
	require.NoError(t, err)

	require.Equal(t, meta, meta2)
}

func TestUserMetaMarshalJSON(t *testing.T) {
	meta := `{"foo":"bar"}`

	data, err := MarshalJSON(meta)
	require.NoError(t, err)

	meta2, err := UnmarshalJSON(data)
	require.NoError(t, err)

	require.Equal(t, meta, meta2)
}

func TestDeepUserMeta(t *testing.T) {
	deepMeta := DeepUserMeta{
		"s": "text",
		"o": map[string]interface{}{
			"a": []interface{}{1.0, 2.0, 3.0},
		},
	}

	meta, err := deepMeta.toUserMeta()
	require.NoError(t, err)
	require.Equal(t, UserMeta{
		"s":      "text",
		"json:o": `{"a":[1,2,3]}`,
	}, meta)

	deepMeta2, err := meta.toDeepUserMeta()
	require.NoError(t, err)
	require.Equal(t, deepMeta, deepMeta2)
}
