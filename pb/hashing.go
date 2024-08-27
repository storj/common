// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.package pb

package pb

import (
	"hash"

	"github.com/zeebo/xxh3"

	"github.com/zeebo/blake3"

	"storj.io/common/pkcrypto"
)

// NewHashFromAlgorithm creates the hash function based on hash algorithm.
func NewHashFromAlgorithm(algorithm PieceHashAlgorithm) hash.Hash {
	switch algorithm {
	case PieceHashAlgorithm_XXH:
		return xxh3.New()
	case PieceHashAlgorithm_BLAKE3:
		return blake3.New()
	default:
		return pkcrypto.NewHash()
	}
}
