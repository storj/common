// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package pb_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
)

func TestMarshalEmptyIsNotNil(t *testing.T) {
	data, err := pb.Marshal(&pb.AuditHistory{})
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Empty(t, data)
}

func TestUnmarshalResets(t *testing.T) {
	raw, err := pb.Marshal(&pb.SignedNodeTagSets{Tags: []*pb.SignedNodeTagSet{{SignerNodeId: []byte{1}}}})
	require.NoError(t, err)

	msg := &pb.SignedNodeTagSets{Tags: []*pb.SignedNodeTagSet{{SignerNodeId: []byte{2}}}}
	require.NoError(t, pb.Unmarshal(raw, msg))
	require.Len(t, msg.Tags, 1)
	require.Equal(t, []byte{1}, msg.Tags[0].SignerNodeId)
}
