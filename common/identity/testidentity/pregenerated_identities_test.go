// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testidentity

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/common/peertls"
	"storj.io/common/storj"
)

func TestPregeneratedIdentity(t *testing.T) {
	IdentityVersionsTest(t, func(t *testing.T, version storj.IDVersion, ident *identity.FullIdentity) {
		assert.Equal(t, version.Number, ident.ID.Version().Number)

		caVersion, err := storj.IDVersionFromCert(ident.CA)
		require.NoError(t, err)
		assert.Equal(t, version.Number, caVersion.Number)

		chains := identity.ToChains(ident.Chain())
		err = peertls.VerifyPeerCertChains(nil, chains)
		assert.NoError(t, err)
	})
}

func TestPregeneratedSignedIdentity(t *testing.T) {
	SignedIdentityVersionsTest(t, func(t *testing.T, version storj.IDVersion, ident *identity.FullIdentity) {
		assert.Equal(t, version.Number, ident.ID.Version().Number)

		caVersion, err := storj.IDVersionFromCert(ident.CA)
		require.NoError(t, err)
		assert.Equal(t, version.Number, caVersion.Number)

		chains := identity.ToChains(ident.Chain())
		err = peertls.VerifyPeerCertChains(nil, chains)
		assert.NoError(t, err)

		signer := NewPregeneratedSigner(ident.ID.Version())
		err = peertls.VerifyCAWhitelist([]*x509.Certificate{signer.Cert})(nil, chains)
		assert.NoError(t, err)
	})
}
