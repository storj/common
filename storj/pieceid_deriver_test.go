// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/identity/testidentity"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestPieceID_PieceDeriver(t *testing.T) {
	pieceIDA := storj.NewPieceID()
	pieceIDB := storj.NewPieceID()
	a := pieceIDA.Deriver()
	b := pieceIDB.Deriver()

	n0 := testidentity.MustPregeneratedIdentity(0, storj.LatestIDVersion()).ID
	n1 := testidentity.MustPregeneratedIdentity(1, storj.LatestIDVersion()).ID

	require.Equal(t, pieceIDA.Derive(n0, 1), a.Derive(n0, 1))
	require.Equal(t, pieceIDB.Derive(n1, 1), b.Derive(n1, 1))

	assert.NotEqual(t, a.Derive(n0, 0), a.Derive(n1, 0), "a(n0, 0) != a(n1, 0)")
	assert.NotEqual(t, b.Derive(n0, 0), b.Derive(n1, 0), "b(n0, 0) != b(n1, 0)")
	assert.NotEqual(t, a.Derive(n0, 0), b.Derive(n0, 0), "a(n0, 0) != b(n0, 0)")
	assert.NotEqual(t, a.Derive(n1, 0), b.Derive(n1, 0), "a(n1, 0) != b(n1, 0)")

	assert.NotEqual(t, a.Derive(n0, 0), a.Derive(n0, 1), "a(n0, 0) != a(n0, 1)")

	// idempotent
	assert.Equal(t, a.Derive(n0, 0), a.Derive(n0, 0), "a(n0, 0)")
	assert.Equal(t, a.Derive(n1, 0), a.Derive(n1, 0), "a(n1, 0)")

	assert.Equal(t, b.Derive(n0, 0), b.Derive(n0, 0), "b(n0, 0)")
	assert.Equal(t, b.Derive(n1, 0), b.Derive(n1, 0), "b(n1, 0)")
}

func BenchmarkDeriver(b *testing.B) {
	pieceID := storj.NewPieceID()
	n0 := testidentity.MustPregeneratedIdentity(0, storj.LatestIDVersion()).ID

	b.Run("Derive", func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			sink = pieceID.Derive(n0, 0)
		}
	})

	b.Run("Deriver", func(b *testing.B) {
		deriver := pieceID.Deriver()
		for k := 0; k < b.N; k++ {
			sink = deriver.Derive(n0, 0)
		}
	})
}

var sink storj.PieceID

func BenchmarkDeriver100(b *testing.B) {
	nodes := make([]storj.NodeID, 100)
	for i := range nodes {
		nodes[i] = testrand.NodeID()
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pieceID := storj.PieceID{byte(i), byte(i >> 8)}
		deriver := pieceID.Deriver()
		for i := range nodes {
			sink = deriver.Derive(nodes[i], int32(i))
		}
	}
}
