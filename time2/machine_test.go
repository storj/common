// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information

package time2_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/time2"
)

func TestMachineWithTimeAt(t *testing.T) {
	expected := time.Now().Add(time.Hour)
	tm := time2.NewMachine(time2.WithTimeAt(expected))
	require.Equal(t, expected, tm.Now())
}

func TestMachineBlockReturnsTrueWhenBlocked(t *testing.T) {
	tm := time2.NewMachine()
	tm.Clock().NewTimer(time.Second)
	require.True(t, tm.Block(context.Background(), 1))
}

func TestMachineBlockRespondsToContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tm := time2.NewMachine()
	require.False(t, tm.Block(ctx, 1))
}
