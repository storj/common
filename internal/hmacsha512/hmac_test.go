// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package hmacsha512_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/internal/hmacsha512"
)

// NodeID is a duplicate of storj.NodeID to avoid circular dependency.
type NodeID [32]byte

// PieceID is a duplicate of storj.PieceID to avoid circular dependency.
type PieceID [32]byte

var sinkSum [hmacsha512.Size]byte
var sinkBytes []byte

func TestRandom(t *testing.T) {
	for i := range 10 {
		key := make([]byte, 21*(i+1))
		node1 := NodeID{}
		node2 := NodeID{}
		_, _ = rand.Read(key)
		_, _ = rand.Read(node1[:])
		_, _ = rand.Read(node2[:])

		var opt hmacsha512.Partial
		opt.Init(key)

		std := hmac.New(sha512.New, key)

		opt.Write(node1[:])
		opt.Write([]byte{1, 0, 0, 0})
		got := opt.SumAndReset()

		std.Reset()
		_, _ = std.Write(node1[:])
		_, _ = std.Write([]byte{1, 0, 0, 0})
		exp := std.Sum(nil)

		require.Equal(t, exp, got[:])

		opt.Write(node1[:])
		opt.Write([]byte{1, 0, 0, 0})
		got = opt.SumAndReset()
		require.Equal(t, exp, got[:])

		opt.Write(node2[:])
		opt.Write([]byte{2, 0, 0, 0})
		got = opt.SumAndReset()

		std.Reset()
		_, _ = std.Write(node2[:])
		_, _ = std.Write([]byte{2, 0, 0, 0})
		exp = std.Sum(nil)

		require.Equal(t, exp, got[:])
	}
}

func BenchmarkInlined(b *testing.B) {
	for i := 0 + 1; i < b.N+1; i++ {
		pieceID := PieceID{byte(i), byte(i), byte(i), byte(i)}

		var mac hmacsha512.Partial
		mac.Init(pieceID[:])
		for k := range 100 {
			nodeid := NodeID{byte(k), byte(k), byte(k), byte(k)}
			mac.Write(nodeid[:])
			var num [4]byte
			binary.BigEndian.PutUint32(num[:], uint32(k))
			mac.Write(num[:])
			sinkSum = mac.SumAndReset()
		}
	}
}

func BenchmarkStandard(b *testing.B) {
	for i := 0 + 1; i < b.N+1; i++ {
		pieceID := PieceID{byte(i), byte(i), byte(i), byte(i)}
		mac := hmac.New(sha512.New, pieceID[:])
		for k := range 100 {
			nodeid := NodeID{byte(k), byte(k), byte(k), byte(k)}
			mac.Reset()
			_, _ = mac.Write(nodeid[:])
			var num [4]byte
			binary.BigEndian.PutUint32(num[:], uint32(k))
			_, _ = mac.Write(num[:])
			sinkBytes = mac.Sum(nil)
		}
	}
}
