// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package pb_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/storj"
)

func TestCopyNode(t *testing.T) {
	node := &pb.Node{
		Id: storj.NodeID{1},
		Address: &pb.NodeAddress{
			Address: "localhost:1234",
			NoiseInfo: &pb.NoiseInfo{
				Proto:     pb.NoiseProtocol_NOISE_IK_25519_AESGCM_BLAKE2B,
				PublicKey: []byte{1, 2, 3},
			},
		},
	}

	copy := pb.CopyNode(node)
	require.EqualValues(t, node, copy)
}
