// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package bloomfilter_test

import (
	"flag"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/bloomfilter"
	"storj.io/common/memory"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestNoFalsePositive(t *testing.T) {
	const numberOfPieces = 10000
	pieceIDs := generateTestIDs(numberOfPieces)

	for _, ratio := range []float32{0.5, 1, 2} {
		size := int64(numberOfPieces * ratio)
		filter := bloomfilter.NewOptimal(size, 0.1)
		for _, pieceID := range pieceIDs {
			filter.Add(pieceID)
		}
		for _, pieceID := range pieceIDs {
			require.True(t, filter.Contains(pieceID))
		}
	}
}

func TestBytes(t *testing.T) {
	for _, count := range []int64{0, 100, 1000, 10000} {
		filter := bloomfilter.NewOptimal(count, 0.1)
		for i := int64(0); i < count; i++ {
			id := testrand.PieceID()
			filter.Add(id)
		}

		bytes := filter.Bytes()
		unmarshaled, err := bloomfilter.NewFromBytes(bytes)
		require.NoError(t, err)

		require.Equal(t, filter, unmarshaled)
	}
}

func TestBytes_Failing(t *testing.T) {
	failing := [][]byte{
		{},
		{0},
		{1},
		{1, 0},
		{255, 10, 10, 10},
	}
	for _, bytes := range failing {
		_, err := bloomfilter.NewFromBytes(bytes)
		require.Error(t, err)
	}
}

// generateTestIDs generates n piece ids.
func generateTestIDs(n int) []storj.PieceID {
	ids := make([]storj.PieceID, n)
	for i := range ids {
		ids[i] = testrand.PieceID()
	}
	return ids
}

func BenchmarkFilterAdd(b *testing.B) {
	ids := generateTestIDs(100000)
	filter := bloomfilter.NewOptimal(int64(len(ids)), 0.1)

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.Add(ids[i%len(ids)])
		}
	})

	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.Contains(ids[i%len(ids)])
		}
	})
}

var approximateFalsePositives = flag.Bool("approximate-false-positive-rate", false, "")

func TestApproximateFalsePositives(t *testing.T) {
	if !*approximateFalsePositives {
		t.Skip("Use --approximate-false-positive-rate to enable diagnostic test.")
	}

	const measurements = 100
	const validation = 1000

	for _, p := range []float64{0.01, 0.05, 0.1, 0.2, 0.3} {
		for _, n := range []int64{1000, 10000, 100000, 1000000} {
			fpp := []float64{}

			for k := 0; k < measurements; k++ {
				filter := bloomfilter.NewOptimal(n, p)
				for i := int64(0); i < n; i++ {
					filter.Add(testrand.PieceID())
				}

				positive := 0
				for k := 0; k < validation; k++ {
					if filter.Contains(testrand.PieceID()) {
						positive++
					}
				}
				fpp = append(fpp, float64(positive)/validation)
			}

			hashCount, size := bloomfilter.NewOptimal(n, p).Parameters()
			summary := summarize(p, fpp)
			t.Logf("n=%8d p=%.2f avg=%.2f min=%.2f mean=%.2f max=%.2f mse=%.3f hc=%d sz=%s", n, p, summary.avg, summary.min, summary.mean, summary.max, summary.mse, hashCount, memory.Size(size))
		}
	}
}

func TestAddFilter(t *testing.T) {
	doesNotContainAtLeastOne := func(filter *bloomfilter.Filter, ids []storj.PieceID) bool {
		for _, id := range ids {
			if !filter.Contains(id) {
				return true
			}
		}
		return false
	}

	ids1 := generateTestIDs(50000)
	ids2 := generateTestIDs(50000)

	filter1 := bloomfilter.NewOptimal(25000, 0.1)
	for _, id := range ids1 {
		filter1.Add(id)
	}

	filter2 := bloomfilter.NewExplicit(filter1.SeedAndParameters())
	for _, id := range ids2 {
		filter2.Add(id)
	}

	require.True(t, doesNotContainAtLeastOne(filter1, ids2), "at least one ID from the 2nd set should not be contained before merge")

	err := filter1.AddFilter(filter2)
	require.NoError(t, err)

	require.False(t, doesNotContainAtLeastOne(filter1, ids2), "all IDs from the 2nd set should be contained after merge")
}

func TestAddFilter_Bad(t *testing.T) {
	t.Run("mismatched seed", func(t *testing.T) {
		filter1 := bloomfilter.NewExplicit(100, 4, 300)
		filter2 := bloomfilter.NewExplicit(101, 4, 300)
		err := filter1.AddFilter(filter2)
		require.EqualError(t, err, "cannot merge: mismatched seed: expected 100 but got 101")
	})
	t.Run("mismatched heap count", func(t *testing.T) {
		filter1 := bloomfilter.NewExplicit(100, 4, 300)
		filter2 := bloomfilter.NewExplicit(100, 5, 300)
		err := filter1.AddFilter(filter2)
		require.EqualError(t, err, "cannot merge: mismatched hash count: expected 4 but got 5")
	})
	t.Run("mismatched table size", func(t *testing.T) {
		filter1 := bloomfilter.NewExplicit(100, 4, 300)
		filter2 := bloomfilter.NewExplicit(100, 4, 400)
		err := filter1.AddFilter(filter2)
		require.EqualError(t, err, "cannot merge: mismatched table size: expected 300 but got 400")
	})
}

type stats struct {
	avg, min, mean, max, mse float64
}

// summarize calculates average, minimum, maximum and mean squared error.
func summarize(expected float64, values []float64) (r stats) {
	sort.Float64s(values)

	for _, v := range values {
		r.avg += v
		r.mse += (v - expected) * (v - expected)
	}
	r.avg /= float64(len(values))
	r.mse /= float64(len(values))

	r.min, r.mean, r.max = values[0], values[len(values)/2], values[len(values)-1]

	return r
}
