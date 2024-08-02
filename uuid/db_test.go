// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package uuid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestNullUUID_SpannerEncoding(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		original := uuid.NullUUID{
			UUID:  uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
			Valid: true,
		}

		encoded, err := original.EncodeSpanner()
		require.NoError(t, err)
		var res uuid.NullUUID
		err = res.DecodeSpanner(encoded)
		require.NoError(t, err)
		require.Equal(t, original, res)
	})

	t.Run("invalid", func(t *testing.T) {
		original := uuid.NullUUID{
			UUID:  uuid.UUID{},
			Valid: false,
		}

		encoded, err := original.EncodeSpanner()
		require.NoError(t, err)
		var res uuid.NullUUID
		err = res.DecodeSpanner(encoded)
		require.NoError(t, err)
		require.Equal(t, original, res)
	})

	t.Run("null bytes decoding", func(t *testing.T) {
		// a NULL BYTES column is returned as an empty uninitialized byte slice
		tests := []struct {
			name    string
			input   any
			want    uuid.NullUUID
			wantErr bool
		}{
			{
				name:  "nil succeeds and is invalid",
				input: nil,
				want: uuid.NullUUID{
					UUID:  uuid.UUID{},
					Valid: false,
				},
				wantErr: false,
			},
			{
				name:  "empty instantiated byte slice fails and errors",
				input: []byte{},
				want: uuid.NullUUID{
					UUID:  uuid.UUID{},
					Valid: true,
				},
				wantErr: true,
			},
			{
				name:  "instantiated byte slice with nil succeeds and errors",
				input: []byte{},
				want: uuid.NullUUID{
					UUID:  uuid.UUID{},
					Valid: true,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var n uuid.NullUUID
				err := n.DecodeSpanner(tt.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				}
				assert.Equal(t, tt.want, n)
			})
		}
	})
}
