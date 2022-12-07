// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information

package time2_test

import (
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
