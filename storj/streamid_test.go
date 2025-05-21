// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestStreamID_Encode(t *testing.T) {
	for range 10 {
		expectedSize := testrand.Intn(255)
		streamID := testrand.StreamID(expectedSize)

		fromString, err := storj.StreamIDFromString(streamID.String())
		require.NoError(t, err)
		require.Equal(t, streamID.String(), fromString.String())

		fromBytes, err := storj.StreamIDFromBytes(streamID.Bytes())
		require.NoError(t, err)
		require.Equal(t, streamID.Bytes(), fromBytes.Bytes())

		require.Equal(t, streamID, fromString)
		require.Equal(t, expectedSize, fromString.Size())
		require.Equal(t, streamID, fromBytes)
		require.Equal(t, expectedSize, fromBytes.Size())
	}
}
