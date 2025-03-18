// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/memory"
	"storj.io/common/storj"
	"storj.io/common/testcontext"
)

func TestRedundancySchemeStripesCount(t *testing.T) {
	scheme := storj.RedundancyScheme{
		ShareSize:      1,
		RequiredShares: 8,
	}

	cases := []struct {
		EncryptedSize int32
		StripesLen    int32
	}{
		{
			EncryptedSize: 1,
			StripesLen:    1,
		},
		{
			EncryptedSize: 7,
			StripesLen:    1,
		},
		{
			EncryptedSize: 57,
			StripesLen:    8,
		},
		{
			EncryptedSize: 63,
			StripesLen:    8,
		},
		{
			EncryptedSize: 64,
			StripesLen:    8,
		},
		{
			EncryptedSize: 65,
			StripesLen:    9,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.StripesLen, scheme.StripeCount(c.EncryptedSize))
	}
}

func TestRedundancyPieceSize(t *testing.T) {
	const uint32Size = 4

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	type TestCase struct {
		Size         int64
		ExpectedSize int64
	}

	// ExpectedSize was precalcualted to avoid dependency to uplink/private/eestream
	// and vivint/infectious using code:
	//    func calcPieceSize(size int64) int {
	//       fc, _ := infectious.NewFEC(2, 4)
	//
	// 	     es := eestream.NewRSScheme(fc, 1*memory.KiB.Int())
	// 	     rs, _ := eestream.NewRedundancyStrategy(es, 0, 0)
	//
	// 	     randReader := io.NopCloser(io.LimitReader(testrand.Reader(), size))
	// 	     readers, _ := eestream.EncodeReader2(context.Background(), encryption.PadReader(randReader, es.StripeSize()), rs)
	//
	// 	     piece, _ := io.ReadAll(readers[0])
	// 	     return len(piece)
	//    }

	for i, tc := range []TestCase{
		{0, 1024},
		{1, 1024},
		{1*memory.KiB.Int64() - uint32Size, 1024},
		{1 * memory.KiB.Int64(), 1024},
		{32*memory.KiB.Int64() - uint32Size, 16384},
		{32 * memory.KiB.Int64(), 17408},
		{32*memory.KiB.Int64() + 100, 17408},
	} {
		errTag := fmt.Sprintf("%d. %+v", i, tc.Size)

		redundancy := storj.RedundancyScheme{
			RequiredShares: 2,
			TotalShares:    4,
			ShareSize:      1 * memory.KiB.Int32(),
		}

		require.Equal(t, tc.ExpectedSize, redundancy.PieceSize(tc.Size), errTag)
	}
}

func TestRedundancyScheme_DB_EncodeDecode(t *testing.T) {
	schemeIn := storj.RedundancyScheme{
		Algorithm:      storj.ReedSolomon,
		ShareSize:      1,
		RepairShares:   2,
		RequiredShares: 3,
		OptimalShares:  4,
		TotalShares:    5,
	}

	value, err := schemeIn.Value()
	require.NoError(t, err)

	require.Equal(t, int64(361416082004640001), value)

	var schemOut storj.RedundancyScheme

	err = schemOut.Scan(int64(0))
	require.NoError(t, err)
	require.Equal(t, storj.RedundancyScheme{}, schemOut)

	err = schemOut.Scan(value)
	require.NoError(t, err)

	require.Equal(t, schemeIn, schemOut)

	valueSpanner, err := schemeIn.EncodeSpanner()
	require.NoError(t, err)

	var schemOutSpanner storj.RedundancyScheme
	err = schemOutSpanner.DecodeSpanner(valueSpanner)
	require.NoError(t, err)

	require.Equal(t, schemeIn, schemOutSpanner)

	err = schemOut.Scan("invalid_value")
	require.Error(t, err)

	err = schemOutSpanner.Scan([]byte("invalid_value_spanner"))
	require.Error(t, err)

	for _, wrongSchem := range []storj.RedundancyScheme{
		{ShareSize: -1},
		{RequiredShares: -1},
		{RepairShares: -1},
		{OptimalShares: -1},
		{TotalShares: -1},
	} {
		_, err := wrongSchem.Value()
		require.Error(t, err)
	}
}
