// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package uuid_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/uuid"
)

func TestValuer(t *testing.T) {
	expected := uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	var a, b uuid.UUID
	err := a.Scan("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	require.NoError(t, err)
	require.Equal(t, expected, a)

	err = b.Scan([]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8})
	require.NoError(t, err)
	require.Equal(t, expected, b)
}

func TestNullValuer(t *testing.T) {
	expected := uuid.NullUUID{
		UUID:  uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
		Valid: true,
	}

	var a, b uuid.NullUUID
	err := a.Scan("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	require.NoError(t, err)
	require.Equal(t, expected, a)

	err = a.Scan(nil)
	require.NoError(t, err)
	require.Equal(t, uuid.NullUUID{}, a)

	err = b.Scan(nil)
	require.NoError(t, err)
	require.Equal(t, uuid.NullUUID{}, a)

	err = b.Scan([]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8})
	require.NoError(t, err)
	require.Equal(t, expected, b)
}
