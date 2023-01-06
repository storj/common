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
	var a int32
	wait := sync2.Go(
		func() { atomic.AddInt32(&a, 1) },
		func() { atomic.AddInt32(&a, 1) },
	)
	wait()
	require.Equal(t, int32(2), a)
	wait()
	require.Equal(t, int32(2), a)
}
