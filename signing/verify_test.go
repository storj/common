// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package signing_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/common/pb"
	"storj.io/common/signing"
	"storj.io/common/storj"
)

func BenchmarkVerifyUplinkOrderSignature(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()

	publicKeyBytes, _ := hex.DecodeString("01eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")
	privateKeyBytes, _ := hex.DecodeString("afefcccadb3d17b1f241b7c83f88c088b54c01b5a25409c13cbeca6bfa22b06901eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")

	publicKey, err := storj.PiecePublicKeyFromBytes(publicKeyBytes)
	require.NoError(b, err)
	privateKey, err := storj.PiecePrivateKeyFromBytes(privateKeyBytes)
	require.NoError(b, err)
	_, _ = publicKey, privateKey

	signed := `0a1052fdfc072182654f163f5f0f9a621d7210e8071a4017871739c3d458737bf24bf214a7387552b18ad75afc3636974cb0d768901a85446954d59a291dde7fde0c648a242863891f543121d4633778c5b6057e62e607`
	signedBytes, err := hex.DecodeString(signed)
	require.NoError(b, err)

	order := pb.Order{}
	err = pb.Unmarshal(signedBytes, &order)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = signing.VerifyUplinkOrderSignature(ctx, publicKey, &order)
		require.NoError(b, err)
	}
}

func BenchmarkVerifyUplinkPieceHashSignature(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()

	publicKeyBytes, _ := hex.DecodeString("01eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")

	publicKey, err := storj.PiecePublicKeyFromBytes(publicKeyBytes)
	require.NoError(b, err)

	signed := `0a2052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649122081855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f69991a40757ff5203925e02c246babdd91c9321265a158d19c99258493fe5cb6482d4bbbb97dea35227ba7b693a3c878e47d8392fc78388e225b541b98c799be7fce3c0720e8072a0c08ba92a3e90510c89afe8202`
	signedBytes, err := hex.DecodeString(signed)
	require.NoError(b, err)

	hash := pb.PieceHash{}
	err = pb.Unmarshal(signedBytes, &hash)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = signing.VerifyUplinkPieceHashSignature(ctx, publicKey, &hash)
		require.NoError(b, err)
	}
}

func BenchmarkVerifyOrderLimitSignature(b *testing.B) {
	b.ReportAllocs()
	ctx := context.Background()

	signer, err := identity.FullIdentityFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\nMIIBYjCCAQigAwIBAgIRAMM/5SHfNFMLl9uTAAQEoZAwCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAS/\n9wOAe42DV90jcRJMMeGe9os528RNJbMthDMkAn58KyOH87Rvlz0uCRnhhk3AbDE+\nXXHfEyed/HPFEMxJwmlGoz8wPTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\nKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0EAwIDSAAw\nRQIhALl9VMhM6NFnPblqOsIHOznsKr0OfQREf/+GSk/t8McsAiAxyOYg3IlB9iA0\nq/pD+qUwXuS+NFyVGOhgdNDFT3amOA==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWzCCAQGgAwIBAgIRAMfle+YJvbpRwr+FqiTrRyswCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARL\nO4n2UCp66X/MY5AzhZsfbBYOBw81Dv8V3y1BXXtbHNsUWNY8RT7r5FSTuLHsaXwq\nTwHdU05bjgnLZT/XdwqaozgwNjAOBgNVHQ8BAf8EBAMCAgQwEwYDVR0lBAwwCgYI\nKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiEA2vce\nasP0sjt6QRJNkgdV/IONJCF0IGgmsCoogCbh9ggCIA3mHgivRBId7sSAU4UUPxpB\nOOfce7bVuJlxvsnNfkkz\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWjCCAQCgAwIBAgIQdzcArqh7Yp9aGiiJXM4+8TAKBggqhkjOPQQDAjAQMQ4w\nDAYDVQQKEwVTdG9yajAiGA8wMDAxMDEwMTAwMDAwMFoYDzAwMDEwMTAxMDAwMDAw\nWjAQMQ4wDAYDVQQKEwVTdG9yajBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM/W\nTxYhs/yGKSg8+Hb2Z/NB2KJef+fWkq7mHl7vhD9JgFwVMowMEFtKOCAhZxLBZD47\nxhYDhHBv4vrLLS+m3wGjODA2MA4GA1UdDwEB/wQEAwICBDATBgNVHSUEDDAKBggr\nBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIC+gM/sI\nXXHq5jJmolw50KKVHlqaqpdxjxJ/6x8oqTHWAiEA1w9EbqPXQ5u/oM+ODf1TBkms\nN9NfnJsY1I2A3NKEvq8=\n-----END CERTIFICATE-----\n"),
		[]byte("-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgzsFsVqt/GdqQlIIJ\nHH2VQNndv1A1fTk/35VPNzLW04ehRANCAATzXrIfcBZAHHxPdFD2PFRViRwe6eWf\nQipaF4iXQmHAW79X4mDx0BibjFfvmzurnYSlyIMZn3jp9RzbLMfnA10C\n-----END PRIVATE KEY-----\n"),
	)
	require.NoError(b, err)

	signee := signing.SignerFromFullIdentity(signer)

	signed := `0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e3802420b088092b8c398feffffff014a0c0899eea2e9051090b98af602524630440220751ae9aa91e78cf5fc858419675cb1148886acfd313c4126870d86c938675e2002206bf29b5efe3752a348446d54d3f10273bc1d582b54cbc2341db7e11508e522085a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d620b088092b8c398feffffff016a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dcc`
	signedBytes, err := hex.DecodeString(signed)

	require.NoError(b, err)

	orderLimit := pb.OrderLimit{}
	err = pb.Unmarshal(signedBytes, &orderLimit)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = signing.VerifyOrderLimitSignature(ctx, signee, &orderLimit)
		require.NoError(b, err)
	}
}
