// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/storj"
)

func TestNodeURL(t *testing.T) {
	id, err := storj.NodeIDFromString("12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		type Test struct {
			String   string
			Expected storj.NodeURL
		}

		for _, testcase := range []Test{
			{"", storj.NodeURL{}},
			// host
			{"33.20.0.1:7777", storj.NodeURL{Address: "33.20.0.1:7777"}},
			{"[2001:db8:1f70::999:de8:7648:6e8]:7777", storj.NodeURL{Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"}},
			{"example.com:7777", storj.NodeURL{Address: "example.com:7777"}},
			// node id + host
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777", storj.NodeURL{ID: id, Address: "33.20.0.1:7777"}},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777", storj.NodeURL{ID: id, Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"}},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com:7777", storj.NodeURL{ID: id, Address: "example.com:7777"}},
			// node id
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@", storj.NodeURL{ID: id}},
			// debounce_limit
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@?debounce=3", storj.NodeURL{ID: id, DebounceLimit: 3}},
			// features
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@?f=ff", storj.NodeURL{ID: id, Features: 255}},
			// noise
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777?noise_proto=2",
				storj.NodeURL{
					ID:      id,
					Address: "33.20.0.1:7777",
					NoiseInfo: storj.NoiseInfo{
						Proto: storj.NoiseProto_IK_25519_AESGCM_BLAKE2b,
					},
				},
			},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777?noise_pub=12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7",
				storj.NodeURL{
					ID:      id,
					Address: "33.20.0.1:7777",
					NoiseInfo: storj.NoiseInfo{
						PublicKey: string(id.Bytes()),
					},
				},
			},
		} {
			url, err := storj.ParseNodeURL(testcase.String)
			require.NoError(t, err, testcase.String)

			assert.Equal(t, testcase.Expected, url)
			assert.Equal(t, testcase.String, url.String())

			copy := pb.NodeFromNodeURL(url).NodeURL()
			assert.Equal(t, testcase.Expected, copy)
			assert.Equal(t, testcase.String, copy.String())
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for _, testcase := range []string{
			// invalid host
			"exampl e.com:7777",
			// invalid node id
			"12vha9oTFnerxgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777",
			"12vha9oTFnerx YRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777",
			"12vha9oTFnerxYRgeQ2BZqoFrLrn_5UWTCY2jA77dF3YvWew7@example.com:7777",
			// invalid node id
			"1112vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@",
		} {
			_, err := storj.ParseNodeURL(testcase)
			assert.Error(t, err, testcase)
		}
	})
}

func TestNodeURLs(t *testing.T) {
	id, err := storj.NodeIDFromString("12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")
	require.NoError(t, err)

	s := "33.20.0.1:7777," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com?noise_proto=2," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com?noise_proto=1," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@"
	urls, err := storj.ParseNodeURLs(s)
	require.NoError(t, err)
	require.Equal(t, storj.NodeURLs{
		storj.NodeURL{Address: "33.20.0.1:7777"},
		storj.NodeURL{ID: id, Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"},
		storj.NodeURL{ID: id, Address: "example.com"},
		storj.NodeURL{ID: id, Address: "example.com", NoiseInfo: storj.NoiseInfo{Proto: storj.NoiseProto_IK_25519_AESGCM_BLAKE2b}},
		pb.NodeFromNodeURL(storj.NodeURL{ID: id, Address: "example.com", NoiseInfo: storj.NoiseInfo{Proto: storj.NoiseProto_IK_25519_ChaChaPoly_BLAKE2b}}).NodeURL(),
		storj.NodeURL{ID: id},
	}, urls)

	require.Equal(t, s, urls.String())
}
