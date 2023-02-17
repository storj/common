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

func TestRedundancyScheme_DownloadNodes(t *testing.T) {
	for i, tt := range []struct {
		k, m, o, n int16
		needed     int32
	}{
		{k: 0, m: 0, o: 0, n: 0, needed: 0},
		{k: 1, m: 1, o: 1, n: 1, needed: 1},
		{k: 1, m: 1, o: 2, n: 2, needed: 2},
		{k: 1, m: 2, o: 2, n: 2, needed: 2},
		{k: 2, m: 3, o: 4, n: 4, needed: 3},
		{k: 2, m: 4, o: 6, n: 8, needed: 3},
		{k: 20, m: 30, o: 40, n: 50, needed: 25},
		{k: 29, m: 35, o: 80, n: 95, needed: 34},
	} {
		tag := fmt.Sprintf("#%d. %+v", i, tt)

		rs := storj.RedundancyScheme{
			RequiredShares: tt.k,
			RepairShares:   tt.m,
			OptimalShares:  tt.o,
			TotalShares:    tt.n,
		}

		assert.Equal(t, tt.needed, rs.DownloadNodes(), tag)
	}
}

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
