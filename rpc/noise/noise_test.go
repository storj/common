// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package noise

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/common/pb"
	"storj.io/common/storj"
	"storj.io/common/testcontext"
)

func TestProtoConversion(t *testing.T) {
	for _, proto := range []pb.NoiseProtocol{
		pb.NoiseProtocol_NOISE_IK_25519_CHACHAPOLY_BLAKE2B,
		pb.NoiseProtocol_NOISE_IK_25519_AESGCM_BLAKE2B} {
		cfg, err := ProtoToConfig(proto)
		require.NoError(t, err)
		proto2, err := ConfigToProto(cfg)
		require.NoError(t, err)
		require.Equal(t, proto, proto2)
	}
}

func TestKeyAttestation(t *testing.T) {
	ctx := testcontext.New(t)
	ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)
	ident2, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)
	ident3, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)

	noiseCfg, err := GenerateServerConf(DefaultProto, ident)
	require.NoError(t, err)

	info, err := ConfigToInfo(noiseCfg)
	require.NoError(t, err)

	attestation, err := GenerateKeyAttestation(ctx, ident, info)
	require.NoError(t, err)

	require.NoError(t, ValidateKeyAttestation(ctx, attestation, ident.ID))

	otherID, err := storj.NodeIDFromString("121RTSDpyNZVcEU84Ticf2L1ntiuUimbWgfATz21tuvgk3vzoA6")
	require.NoError(t, err)
	err = ValidateKeyAttestation(ctx, attestation, otherID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "node id mismatch")

	badAttestation2 := *attestation
	badAttestation2.NoisePublicKey = badAttestation2.NoisePublicKey[:len(badAttestation2.NoisePublicKey)-1]
	err = ValidateKeyAttestation(ctx, &badAttestation2, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "signature is not valid")

	badAttestation3 := *attestation
	badAttestation3.Timestamp = time.Now()
	err = ValidateKeyAttestation(ctx, &badAttestation3, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "signature is not valid")

	ident2.CA = ident.CA
	badAttestation4 := *attestation
	badAttestation4.NodeCertchain = identity.EncodePeerIdentity(ident2.PeerIdentity())
	err = ValidateKeyAttestation(ctx, &badAttestation4, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certificate chain invalid")

	ident3.Leaf = ident.Leaf
	badAttestation5 := *attestation
	badAttestation5.NodeCertchain = identity.EncodePeerIdentity(ident3.PeerIdentity())
	err = ValidateKeyAttestation(ctx, &badAttestation5, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certificate chain invalid")
}

func TestNoiseSessionAttestation(t *testing.T) {
	ctx := testcontext.New(t)
	ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)
	ident2, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)
	ident3, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{})
	require.NoError(t, err)

	var hash [32]byte
	_, err = rand.Read(hash[:])
	require.NoError(t, err)

	attestation, err := GenerateSessionAttestation(ctx, ident, hash[:])
	require.NoError(t, err)

	require.NoError(t, ValidateSessionAttestation(ctx, attestation, ident.ID))

	otherID, err := storj.NodeIDFromString("121RTSDpyNZVcEU84Ticf2L1ntiuUimbWgfATz21tuvgk3vzoA6")
	require.NoError(t, err)
	err = ValidateSessionAttestation(ctx, attestation, otherID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "node id mismatch")

	badAttestation2 := *attestation
	badAttestation2.NoiseHandshakeHash = badAttestation2.NoiseHandshakeHash[:len(badAttestation2.NoiseHandshakeHash)-1]
	err = ValidateSessionAttestation(ctx, &badAttestation2, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "signature is not valid")

	ident2.CA = ident.CA
	badAttestation3 := *attestation
	badAttestation3.NodeCertchain = identity.EncodePeerIdentity(ident2.PeerIdentity())
	err = ValidateSessionAttestation(ctx, &badAttestation3, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certificate chain invalid")

	ident3.Leaf = ident.Leaf
	badAttestation4 := *attestation
	badAttestation4.NodeCertchain = identity.EncodePeerIdentity(ident3.PeerIdentity())
	err = ValidateSessionAttestation(ctx, &badAttestation4, ident.ID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "certificate chain invalid")
}
