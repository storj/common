// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package pkcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigningAndVerifyingECDSA(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"empty", ""},
		{"single byte", "C"},
		{"longnulls", string(make([]byte, 2000))},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			privKey, err := GeneratePrivateECDSAKey(authECCurve)
			assert.NoError(t, err)
			pubKey, err := PublicKeyFromPrivate(privKey)
			require.NoError(t, err)

			// test signing and verifying a hash of the data
			sig, err := HashAndSign(privKey, []byte(test.data))
			assert.NoError(t, err)
			err = HashAndVerifySignature(pubKey, []byte(test.data), sig)
			assert.NoError(t, err)

			// test signing and verifying the data directly
			sig, err = SignWithoutHashing(privKey, []byte(test.data))
			assert.NoError(t, err)
			err = VerifySignatureWithoutHashing(pubKey, []byte(test.data), sig)
			assert.NoError(t, err)
		})
	}
}

func TestSigningAndVerifyingRSA(t *testing.T) {
	privKey, err := GeneratePrivateRSAKey(StorjRSAKeyBits)
	assert.NoError(t, err)
	pubKey, err := PublicKeyFromPrivate(privKey)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{"empty", ""},
		{"single byte", "C"},
		{"longnulls", string(make([]byte, 2000))},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			// test signing and verifying a hash of the data
			sig, err := HashAndSign(privKey, []byte(test.data))
			assert.NoError(t, err)
			err = HashAndVerifySignature(pubKey, []byte(test.data), sig)
			assert.NoError(t, err)

			// don't test signing and verifying the data directly, as RSA can't
			// handle messages of arbitrary size
		})
	}
}

func TestPublicKeyFromPrivate(t *testing.T) {
	t.Run("RSA", func(t *testing.T) {
		privKey, err := GeneratePrivateRSAKey(StorjRSAKeyBits)
		require.NoError(t, err)

		pubKey, err := PublicKeyFromPrivate(privKey)
		require.NotNil(t, pubKey, "public key cannot be nil")
		require.NoError(t, err)
	})

	t.Run("ECDSA", func(t *testing.T) {
		privKey, err := GeneratePrivateECDSAKey(authECCurve)
		require.NoError(t, err)

		pubKey, err := PublicKeyFromPrivate(privKey)
		require.NotNil(t, pubKey, "public key cannot be nil")
		require.NoError(t, err)
	})

	t.Run("invalid key", func(t *testing.T) {
		_, err := PublicKeyFromPrivate("invalid")
		require.Error(t, err)
		require.True(t, ErrUnsupportedKey.Has(err), "invalid error class")
	})

}
