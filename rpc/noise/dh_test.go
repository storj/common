// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package noise

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/flynn/noise"
	"github.com/stretchr/testify/require"

	"storj.io/common/testrand"
)

func TestDH25519RFC7748(t *testing.T) {
	mustHex := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}

	// RFC 7748, section 5.2.
	scalar := mustHex("a546e36bf0527c9d3b16154b82465edd62144c0ac1fc5a18506a2244ba449ac4")
	point := mustHex("e6db6867583030db3594c1a424b15f7c726624ec26b3353b10a903a6d0ab1c4c")
	out, err := DH25519.DH(scalar, point)
	require.NoError(t, err)
	require.Equal(t, mustHex("c3da55379de9c6908e94ea4df28d084f32eccf03491c71f754b4075577a28552"), out)

	// RFC 7748, section 6.1 (Alice's keypair).
	priv := mustHex("77076d0a7318a57d3c16c17251b26645df4c2f87ebc0992ab177fba51db92c2a")
	key, err := DH25519.GenerateKeypair(bytes.NewReader(priv))
	require.NoError(t, err)
	require.Equal(t, mustHex("8520f0098930a754748b7ddcb43ef75a0dbf3a0d26381af4eba4a98eaa9b4e6a"), key.Public)
}

func TestDH25519MatchesFlynn(t *testing.T) {
	for range 200 {
		seed := testrand.BytesInt(64)

		// same rng input must give the same keypair.
		a, err := DH25519.GenerateKeypair(bytes.NewReader(seed))
		require.NoError(t, err)
		aFlynn, err := noise.DH25519.GenerateKeypair(bytes.NewReader(seed))
		require.NoError(t, err)
		require.Equal(t, aFlynn, a)

		b, err := DH25519.GenerateKeypair(bytes.NewReader(seed[32:]))
		require.NoError(t, err)

		got, err := DH25519.DH(a.Private, b.Public)
		require.NoError(t, err)
		want, err := noise.DH25519.DH(a.Private, b.Public)
		require.NoError(t, err)
		require.Equal(t, want, got)

		other, err := DH25519.DH(b.Private, a.Public)
		require.NoError(t, err)
		require.Equal(t, got, other)
	}
}

func TestDH25519BadInput(t *testing.T) {
	key, err := DH25519.GenerateKeypair(nil)
	require.NoError(t, err)

	for _, tc := range []struct {
		name         string
		priv, pubkey []byte
	}{
		{"low order point", key.Private, make([]byte, 32)},
		{"short private key", key.Private[:31], key.Public},
		{"short public key", key.Private, key.Public[:31]},
	} {
		for _, dh := range []noise.DHFunc{DH25519, noise.DH25519} {
			_, err := dh.DH(tc.priv, tc.pubkey)
			require.Error(t, err, tc.name)
		}
	}
}

// TestDH25519HandshakeInterop checks that a peer using DH25519 can do an IK
// handshake with a peer using noise.DH25519, in both directions.
func TestDH25519HandshakeInterop(t *testing.T) {
	suite := func(dh noise.DHFunc) noise.CipherSuite {
		return noise.NewCipherSuite(dh, noise.CipherChaChaPoly, noise.HashBLAKE2b)
	}
	for _, tc := range []struct{ initiator, responder noise.DHFunc }{
		{noise.DH25519, DH25519},
		{DH25519, noise.DH25519},
		{DH25519, DH25519},
	} {
		responderKey, err := tc.responder.GenerateKeypair(rand.Reader)
		require.NoError(t, err)
		initiatorKey, err := tc.initiator.GenerateKeypair(rand.Reader)
		require.NoError(t, err)

		ini, err := noise.NewHandshakeState(noise.Config{
			CipherSuite: suite(tc.initiator), Pattern: noise.HandshakeIK, Initiator: true,
			StaticKeypair: initiatorKey, PeerStatic: responderKey.Public,
		})
		require.NoError(t, err)
		res, err := noise.NewHandshakeState(noise.Config{
			CipherSuite: suite(tc.responder), Pattern: noise.HandshakeIK,
			StaticKeypair: responderKey,
		})
		require.NoError(t, err)

		// both payloads are encrypted with keys from all four DH results, so
		// if they decrypt, the two sides agree.
		msg1, _, _, err := ini.WriteMessage(nil, []byte("hello"))
		require.NoError(t, err)
		payload, _, _, err := res.ReadMessage(nil, msg1)
		require.NoError(t, err)
		require.Equal(t, []byte("hello"), payload)

		msg2, _, _, err := res.WriteMessage(nil, []byte("world"))
		require.NoError(t, err)
		payload, _, _, err = ini.ReadMessage(nil, msg2)
		require.NoError(t, err)
		require.Equal(t, []byte("world"), payload)
	}
}

func BenchmarkDH25519(b *testing.B) {
	for _, bc := range []struct {
		name string
		dh   noise.DHFunc
	}{{"flynn", noise.DH25519}, {"storj", DH25519}} {
		b.Run(bc.name, func(b *testing.B) {
			k1, err := bc.dh.GenerateKeypair(rand.Reader)
			require.NoError(b, err)
			k2, err := bc.dh.GenerateKeypair(rand.Reader)
			require.NoError(b, err)
			b.Run("GenerateKeypair", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = bc.dh.GenerateKeypair(rand.Reader)
				}
			})
			b.Run("DH", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = bc.dh.DH(k1.Private, k2.Public)
				}
			})
		})
	}
}
