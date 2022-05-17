// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/identity/testidentity"
	"storj.io/common/storj"
)

func TestNewPieceID(t *testing.T) {
	a := storj.NewPieceID()
	assert.NotEmpty(t, a)

	b := storj.NewPieceID()
	assert.NotEqual(t, a, b)
}

func TestPieceID_Encode(t *testing.T) {
	_, err := storj.PieceIDFromString("likn43kilfzd")
	assert.Error(t, err)

	_, err = storj.PieceIDFromBytes([]byte{1, 2, 3, 4, 5})
	assert.Error(t, err)

	for i := 0; i < 10; i++ {
		pieceid := storj.NewPieceID()

		fromString, err := storj.PieceIDFromString(pieceid.String())
		assert.NoError(t, err)
		fromBytes, err := storj.PieceIDFromBytes(pieceid.Bytes())
		assert.NoError(t, err)

		assert.Equal(t, pieceid, fromString)
		assert.Equal(t, pieceid, fromBytes)
	}
}

func TestPieceID_Derive(t *testing.T) {
	a := storj.NewPieceID()
	b := storj.NewPieceID()

	n0 := testidentity.MustPregeneratedIdentity(0, storj.LatestIDVersion()).ID
	n1 := testidentity.MustPregeneratedIdentity(1, storj.LatestIDVersion()).ID

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

func TestPieceID_MarshalJSON(t *testing.T) {
	pieceid := storj.NewPieceID()
	buf, err := json.Marshal(pieceid)
	require.NoError(t, err)
	assert.Equal(t, string(buf), `"`+pieceid.String()+`"`)
}

func TestPieceID_UnmarshalJSON(t *testing.T) {
	originalPieceID := storj.NewPieceID()
	var pieceid storj.PieceID
	err := json.Unmarshal([]byte(`"`+originalPieceID.String()+`"`), &pieceid)
	require.NoError(t, err)
	assert.Equal(t, pieceid.String(), originalPieceID.String())

	assert.Error(t, json.Unmarshal([]byte(`""`+originalPieceID.String()+`""`), &pieceid))
	assert.Error(t, json.Unmarshal([]byte(`{}`), &pieceid))
}

func BenchmarkDeriver(b *testing.B) {
	pieceID := storj.NewPieceID()
	n0 := testidentity.MustPregeneratedIdentity(0, storj.LatestIDVersion()).ID

	b.Run("Derive", func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			_ = pieceID.Derive(n0, 0)
		}
	})

	b.Run("Deriver", func(b *testing.B) {
		deriver := pieceID.Deriver()
		for k := 0; k < b.N; k++ {
			_ = deriver.Derive(n0, 0)
		}
	})
}
