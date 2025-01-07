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
