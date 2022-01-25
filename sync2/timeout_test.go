// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/sync2"
	"storj.io/common/testcontext"
)

func TestWithTimeout(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	workResult := 0
	timeoutResult := 0

	ctx.Go(func() error {
		sync2.WithTimeout(time.Second, func() {
			workResult = 1
		}, func() {
			timeoutResult = 1
		})
		return nil
	})

	ctx.Wait()

	require.Equal(t, workResult, 1)
	require.Equal(t, timeoutResult, 0)
}

func TestWithTimeout_Fail(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	working := make(chan struct{})

	workResult := 0
	timeoutResult := 0
	allResult := 0

	workCompleted := make(chan struct{})
	timeoutCompleted := make(chan struct{})

	ctx.Go(func() error {
		sync2.WithTimeout(0, func() {
			<-working
			workResult = 1
			close(workCompleted)
		}, func() {
			timeoutResult = 1
			close(timeoutCompleted)
		})
		allResult = 1
		return nil
	})

	<-timeoutCompleted

	require.Equal(t, workResult, 0)
	require.Equal(t, timeoutResult, 1)
	require.Equal(t, allResult, 0)

	close(working)
	<-workCompleted

	require.Equal(t, workResult, 1)
	require.Equal(t, timeoutResult, 1)
}
