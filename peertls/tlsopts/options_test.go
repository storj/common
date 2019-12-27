// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package tlsopts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/common/identity/testidentity"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/storj"
)

func TestOptions_DialOption_error_on_empty_ID(t *testing.T) {
	testidentity.CompleteIdentityVersionsTest(t, func(t *testing.T, version storj.IDVersion, ident *identity.FullIdentity) {
		tlsOptions, err := tlsopts.NewOptions(ident, tlsopts.Config{
			PeerIDVersions: "*",
		}, nil)
		require.NoError(t, err)

		dialOption, err := tlsOptions.DialOption(storj.NodeID{})
		assert.Nil(t, dialOption)
		assert.Error(t, err)
	})
}

func TestOptions_DialUnverifiedIDOption(t *testing.T) {
	testidentity.CompleteIdentityVersionsTest(t, func(t *testing.T, version storj.IDVersion, ident *identity.FullIdentity) {
		tlsOptions, err := tlsopts.NewOptions(ident, tlsopts.Config{
			PeerIDVersions: "*",
		}, nil)
		require.NoError(t, err)

		dialOption := tlsOptions.DialUnverifiedIDOption()
		assert.NotNil(t, dialOption)
	})
}
