// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

// RedundancyScheme specifies the parameters and the algorithm for redundancy.
type RedundancyScheme struct {
	// Algorithm determines the algorithm to be used for redundancy.
	Algorithm RedundancyAlgorithm

	// ShareSize is the size in bytes for each erasure shares.
	ShareSize int32

	// RequiredShares is the minimum number of shares required to recover a
	// stripe, reed-solomon k.
	RequiredShares int16
	// RepairShares is the minimum number of safe shares that can remain
	// before a repair is triggered.
	RepairShares int16
	// OptimalShares is the desired total number of shares for a segment.
	OptimalShares int16
	// TotalShares is the number of shares to encode. If it is larger than
	// OptimalShares, slower uploads of the excess shares will be aborted in
	// order to improve performance.
	TotalShares int16
}

// IsZero returns true if no field in the struct is set to non-zero value.
func (scheme RedundancyScheme) IsZero() bool {
	return scheme == (RedundancyScheme{})
}

// StripeSize is the number of bytes for a stripe.
// Stripes are erasure encoded and split into n shares, where we need k to
// reconstruct the stripe. Therefore a stripe size is the erasure share size
// times the required shares, k.
func (scheme RedundancyScheme) StripeSize() int32 {
	return scheme.ShareSize * int32(scheme.RequiredShares)
}

// StripeCount returns segment's total number of stripes based on segment's encrypted size.
func (scheme RedundancyScheme) StripeCount(encryptedSegmentSize int32) int32 {
	stripeSize := scheme.StripeSize()
	return (encryptedSegmentSize + stripeSize - 1) / stripeSize
}

// PieceSize calculates piece size for give size.
func (scheme RedundancyScheme) PieceSize(size int64) int64 {
	const uint32Size = 4
	stripeSize := int64(scheme.StripeSize())
	stripes := (size + uint32Size + stripeSize - 1) / stripeSize

	encodedSize := stripes * int64(scheme.StripeSize())
	pieceSize := encodedSize / int64(scheme.RequiredShares)
	return pieceSize
}

// RedundancyAlgorithm is the algorithm used for redundancy.
type RedundancyAlgorithm byte

// List of supported redundancy algorithms.
const (
	InvalidRedundancyAlgorithm = RedundancyAlgorithm(iota)
	ReedSolomon
)
