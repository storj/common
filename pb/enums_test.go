// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package pb_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/pb/internal/gogopb"
)

// TestEnumMaps verifies the hand-written enum maps against the gogo/protobuf
// generated ones they replace.
func TestEnumMaps(t *testing.T) {
	require.Equal(t, gogopb.PieceAction_name, pb.PieceAction_name)
	require.Equal(t, gogopb.PieceAction_value, pb.PieceAction_value)
	require.Equal(t, gogopb.PieceHashAlgorithm_name, pb.PieceHashAlgorithm_name)
	require.Equal(t, gogopb.PieceHashAlgorithm_value, pb.PieceHashAlgorithm_value)
	require.Equal(t, gogopb.SettlementWithWindowResponse_Status_name, pb.SettlementWithWindowResponse_Status_name)
	require.Equal(t, gogopb.SettlementWithWindowResponse_Status_value, pb.SettlementWithWindowResponse_Status_value)

	for value, name := range pb.PieceAction_name {
		require.Equal(t, name, pb.PieceAction(value).String())
	}
	for value, name := range pb.PieceHashAlgorithm_name {
		require.Equal(t, name, pb.PieceHashAlgorithm(value).String())
	}
	for value, name := range pb.SettlementWithWindowResponse_Status_name {
		require.Equal(t, name, fmt.Sprint(pb.SettlementWithWindowResponse_Status(value)))
	}
}
