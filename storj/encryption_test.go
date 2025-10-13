// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestNewKey(t *testing.T) {
	t.Run("nil humanReadableKey", func(t *testing.T) {
		t.Parallel()

		key, err := storj.NewKey(nil)
		require.NoError(t, err)
		require.True(t, key.IsZero(), "key isn't zero value")
	})

	t.Run("empty humanReadableKey", func(t *testing.T) {
		t.Parallel()

		key, err := storj.NewKey([]byte{})
		require.NoError(t, err)
		require.True(t, key.IsZero(), "key isn't zero value")
	})

	t.Run("humanReadableKey is of KeySize length", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.Bytes(storj.KeySize)

		key, err := storj.NewKey(humanReadableKey)
		require.NoError(t, err)
		require.Equal(t, humanReadableKey, key[:])
	})

	t.Run("humanReadableKey is shorter than KeySize", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(storj.KeySize))

		key, err := storj.NewKey(humanReadableKey)
		require.NoError(t, err)
		require.Equal(t, humanReadableKey, key[:len(humanReadableKey)])
	})

	t.Run("humanReadableKey is larger than KeySize", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(10) + storj.KeySize + 1)

		key, err := storj.NewKey(humanReadableKey)
		require.NoError(t, err)
		assert.Equal(t, humanReadableKey[:storj.KeySize], key[:])
	})

	t.Run("same human readable key produce the same key", func(t *testing.T) {
		t.Parallel()

		humanReadableKey := testrand.BytesInt(testrand.Intn(10) + storj.KeySize + 1)

		key1, err := storj.NewKey(humanReadableKey)
		require.NoError(t, err)

		key2, err := storj.NewKey(humanReadableKey)
		require.NoError(t, err)

		assert.Equal(t, key1, key2, "keys are equal")
	})
}

func TestKey_IsZero(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var key *storj.Key
		require.True(t, key.IsZero())

		wrapperFn := func(key *storj.Key) bool {
			return key.IsZero()
		}
		require.True(t, wrapperFn(nil))
	})

	t.Run("zero", func(t *testing.T) {
		key := &storj.Key{}
		require.True(t, key.IsZero())
	})

	t.Run("no nil/zero", func(t *testing.T) {
		key := &storj.Key{'k'}
		require.False(t, key.IsZero())
	})
}

func TestEncryptedParameters_Scan(t *testing.T) {
	tmp := storj.EncryptionParameters{}
	require.Error(t, tmp.Scan(byte(32)))
	require.Error(t, tmp.Scan(false))
	require.Error(t, tmp.Scan([]byte{}))

	for _, ep := range []storj.EncryptionParameters{
		{CipherSuite: storj.EncUnspecified, BlockSize: 123},
		{CipherSuite: storj.EncNull, BlockSize: 0},
		{CipherSuite: storj.EncAESGCM, BlockSize: 512},
		{CipherSuite: storj.EncSecretBox, BlockSize: 1024},
		{CipherSuite: storj.EncNullBase64URL, BlockSize: 41232},
	} {
		v, err := ep.Value()
		require.NoError(t, err)

		x, ok := v.(int64)
		require.True(t, ok)

		var tmp storj.EncryptionParameters
		err = tmp.Scan(x)
		require.NoError(t, err)
		require.Equal(t, ep, tmp)

		s := strconv.FormatInt(x, 10)
		err = tmp.DecodeSpanner(s)
		require.NoError(t, err)
		require.Equal(t, ep, tmp)
	}
}

// TestNonce_Scan tests (*Nonce).Scan().
func TestNonce_Scan(t *testing.T) {
	tmp := storj.Nonce{}
	require.Error(t, tmp.Scan(32))
	require.Error(t, tmp.Scan(false))
	require.Error(t, tmp.Scan([]byte{}))

	require.NoError(t, tmp.Scan(nil))
	require.True(t, tmp.IsZero())
	require.NoError(t, tmp.Scan(tmp.Bytes()))
	require.True(t, tmp.IsZero())
}

// TestEncryptedPrivateKey_Scan tests (*EncryptedPrivateKey).Scan().
func TestEncryptedPrivateKey_Scan(t *testing.T) {
	tmp := storj.EncryptedPrivateKey{}
	require.Error(t, tmp.Scan(32))
	require.Error(t, tmp.Scan(false))
	require.NoError(t, tmp.Scan([]byte{}))
	require.NoError(t, tmp.Scan([]byte{1, 2, 3, 4}))

	ref := []byte{1, 2, 3}
	require.NoError(t, tmp.Scan(ref))
	ref[0] = 0xFF
	require.Equal(t, storj.EncryptedPrivateKey{1, 2, 3}, tmp)
}
