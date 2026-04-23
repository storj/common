// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
)

func TestGo(t *testing.T) {
	var a atomic.Int32
	wait := sync2.Go(
		func() { a.Add(1) },
		func() { a.Add(1) },
	)
	wait()
	require.Equal(t, int32(2), a.Load())
	wait()
	require.Equal(t, int32(2), a.Load())
}

func TestParallel(t *testing.T) {
	values := []int64{1, 3, 7}

	var total atomic.Int64
	sync2.Parallel(values, func(t int64) {
		total.Add(t)
	})

	require.Equal(t, int64(11), total.Load())
}
