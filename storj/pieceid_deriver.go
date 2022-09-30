// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"hash"
)

// PieceIDDeriver can be used to for multiple derivation from the same PieceID
// without need to initialize mac for each Derive call.
type PieceIDDeriver struct {
	mac hash.Hash
}

// Deriver creates piece ID dervier for multiple derive operations.
func (id PieceID) Deriver() PieceIDDeriver {
	return PieceIDDeriver{
		mac: hmac.New(sha512.New, id.Bytes()),
	}
}

// Derive a new PieceID from the piece ID, the given storage node ID and piece number.
// Initial mac is created from piece ID once while creating PieceDeriver and just
// reset to initial state at the beginning of each call.
func (pd PieceIDDeriver) Derive(storagenodeID NodeID, pieceNum int32) PieceID {
	pd.mac.Reset()

	_, _ = pd.mac.Write(storagenodeID.Bytes()) // on hash.Hash write never returns an error
	num := make([]byte, 4)
	binary.BigEndian.PutUint32(num, uint32(pieceNum))
	_, _ = pd.mac.Write(num) // on hash.Hash write never returns an error
	var derived PieceID
	copy(derived[:], pd.mac.Sum(nil))
	return derived
}
