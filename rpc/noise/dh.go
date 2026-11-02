// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package noise

import (
	"crypto/rand"
	"crypto/subtle"
	"io"

	"filippo.io/edwards25519"
	"filippo.io/edwards25519/field"
	"github.com/flynn/noise"
)

// DH25519 is a noise.DHFunc for Curve25519. It is compatible with
// noise.DH25519 (same keys, same shared secrets), but it is about twice as
// fast: noise.DH25519 goes through crypto/ecdh.NewPrivateKey, which also
// computes a public key that DH does not need.
var DH25519 noise.DHFunc = dh25519{}

type dh25519 struct{}

// GenerateKeypair reads the private key from rng exactly like
// noise.DH25519 does, so the same rng gives the same keypair.
func (dh25519) GenerateKeypair(rng io.Reader) (noise.DHKey, error) {
	if rng == nil {
		rng = rand.Reader
	}
	privkey := make([]byte, 32)
	if _, err := io.ReadFull(rng, privkey); err != nil {
		return noise.DHKey{}, Error.Wrap(err)
	}
	// X25519(k, 9) is the Montgomery form of [k]B on Edwards25519, and the
	// fixed-base multiplication is faster than a Montgomery ladder.
	s, err := new(edwards25519.Scalar).SetBytesWithClamping(privkey)
	if err != nil {
		return noise.DHKey{}, Error.Wrap(err)
	}
	pubkey := new(edwards25519.Point).ScalarBaseMult(s).BytesMontgomery()
	return noise.DHKey{Private: privkey, Public: pubkey}, nil
}

// DH computes the X25519 shared secret. Like noise.DH25519, it returns
// an error for a low order point (all-zero output).
func (dh25519) DH(privkey, pubkey []byte) ([]byte, error) {
	if len(privkey) != 32 {
		return nil, Error.New("bad private key length: %d", len(privkey))
	}
	if len(pubkey) != 32 {
		return nil, Error.New("bad public key length: %d", len(pubkey))
	}
	out := make([]byte, 32)
	x25519ScalarMult(out, privkey, pubkey)
	var zero [32]byte
	if subtle.ConstantTimeCompare(out, zero[:]) == 1 {
		return nil, Error.New("bad input point: low order point")
	}
	return out, nil
}

func (dh25519) DHLen() int     { return 32 }
func (dh25519) DHName() string { return "25519" }

// x25519ScalarMult is the RFC 7748 Montgomery ladder, copied from
// x25519ScalarMult in Go 1.26 crypto/ecdh, which does not export it.
// The caller must check that scalar and point are 32 bytes long.
func x25519ScalarMult(dst, scalar, point []byte) {
	var e [32]byte

	copy(e[:], scalar)
	e[0] &= 248
	e[31] &= 127
	e[31] |= 64

	var x1, x2, z2, x3, z3, tmp0, tmp1 field.Element
	_, _ = x1.SetBytes(point)
	x2.One()
	x3.Set(&x1)
	z3.One()

	swap := 0
	for pos := 254; pos >= 0; pos-- {
		b := e[pos/8] >> uint(pos&7)
		b &= 1
		swap ^= int(b)
		x2.Swap(&x3, swap)
		z2.Swap(&z3, swap)
		swap = int(b)

		tmp0.Subtract(&x3, &z3)
		tmp1.Subtract(&x2, &z2)
		x2.Add(&x2, &z2)
		z2.Add(&x3, &z3)
		z3.Multiply(&tmp0, &x2)
		z2.Multiply(&z2, &tmp1)
		tmp0.Square(&tmp1)
		tmp1.Square(&x2)
		x3.Add(&z3, &z2)
		z2.Subtract(&z3, &z2)
		x2.Multiply(&tmp1, &tmp0)
		tmp1.Subtract(&tmp1, &tmp0)
		z2.Square(&z2)

		z3.Mult32(&tmp1, 121666)
		x3.Square(&x3)
		tmp0.Add(&tmp0, &z3)
		z3.Multiply(&x1, &z2)
		z2.Multiply(&tmp1, &tmp0)
	}

	x2.Swap(&x3, swap)
	z2.Swap(&z3, swap)

	z2.Invert(&z2)
	x2.Multiply(&x2, &z2)
	copy(dst, x2.Bytes())
}
