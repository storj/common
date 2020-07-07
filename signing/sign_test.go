// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package signing_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/common/identity/testidentity"
	"storj.io/common/pb"
	"storj.io/common/signing"
	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
)

// printNewSigned, use when you need to generate a valid signature for tests.
const printNewSigned = false

func TestOrderLimitVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	signer, err := identity.FullIdentityFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\nMIIBYjCCAQigAwIBAgIRAMM/5SHfNFMLl9uTAAQEoZAwCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAS/\n9wOAe42DV90jcRJMMeGe9os528RNJbMthDMkAn58KyOH87Rvlz0uCRnhhk3AbDE+\nXXHfEyed/HPFEMxJwmlGoz8wPTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\nKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0EAwIDSAAw\nRQIhALl9VMhM6NFnPblqOsIHOznsKr0OfQREf/+GSk/t8McsAiAxyOYg3IlB9iA0\nq/pD+qUwXuS+NFyVGOhgdNDFT3amOA==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWzCCAQGgAwIBAgIRAMfle+YJvbpRwr+FqiTrRyswCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARL\nO4n2UCp66X/MY5AzhZsfbBYOBw81Dv8V3y1BXXtbHNsUWNY8RT7r5FSTuLHsaXwq\nTwHdU05bjgnLZT/XdwqaozgwNjAOBgNVHQ8BAf8EBAMCAgQwEwYDVR0lBAwwCgYI\nKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiEA2vce\nasP0sjt6QRJNkgdV/IONJCF0IGgmsCoogCbh9ggCIA3mHgivRBId7sSAU4UUPxpB\nOOfce7bVuJlxvsnNfkkz\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWjCCAQCgAwIBAgIQdzcArqh7Yp9aGiiJXM4+8TAKBggqhkjOPQQDAjAQMQ4w\nDAYDVQQKEwVTdG9yajAiGA8wMDAxMDEwMTAwMDAwMFoYDzAwMDEwMTAxMDAwMDAw\nWjAQMQ4wDAYDVQQKEwVTdG9yajBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM/W\nTxYhs/yGKSg8+Hb2Z/NB2KJef+fWkq7mHl7vhD9JgFwVMowMEFtKOCAhZxLBZD47\nxhYDhHBv4vrLLS+m3wGjODA2MA4GA1UdDwEB/wQEAwICBDATBgNVHSUEDDAKBggr\nBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIC+gM/sI\nXXHq5jJmolw50KKVHlqaqpdxjxJ/6x8oqTHWAiEA1w9EbqPXQ5u/oM+ODf1TBkms\nN9NfnJsY1I2A3NKEvq8=\n-----END CERTIFICATE-----\n"),
		[]byte("-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgzsFsVqt/GdqQlIIJ\nHH2VQNndv1A1fTk/35VPNzLW04ehRANCAATzXrIfcBZAHHxPdFD2PFRViRwe6eWf\nQipaF4iXQmHAW79X4mDx0BibjFfvmzurnYSlyIMZn3jp9RzbLMfnA10C\n-----END PRIVATE KEY-----\n"),
	)
	require.NoError(t, err)

	signee := signing.SignerFromFullIdentity(signer)

	type Hex struct {
		Unsigned string
		Signed   string
	}

	hexes := []Hex{
		{ // 385c0467
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e3802420c08a1dba2e90510e0a1b3c6014a0c08a1dba2e90510e0a1b3c6015a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e3802420c08a1dba2e90510e0a1b3c6014a0c08a1dba2e90510e0a1b3c60152473045022100ada5fc332dfbd607216e961bede421e43e2e336acab8eab2244f5e3a696ede720220365e78e738c19fc9d3cb26b061dcf6439ea702cb0ef1408cf7aeb27cabee4cc45a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d",
		},
		{ // 385c0467 without satellite address
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e3802420c08a1dba2e90510e0a1b3c6014a0c08a1dba2e90510e0a1b3c601",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e3802420c08a1dba2e90510e0a1b3c6014a0c08a1dba2e90510e0a1b3c60152473045022100a2e7849a4cb93e6bbc591949f93a5e97d9b1392a5770667afc634389355e094102200bfa72531afc9359181f7fc5181e387a03dea8234f74a7d7c44ca3aa0c5ab21d",
		},
		{ // 385c0467 without piece expiration
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e38024a0c08a1dba2e90510e0a1b3c6015a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac30001a209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679390022206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746002a2095af5a25367951baa2ff6cd471c483f15fb90badb37c5821b6d95526a41a950430904e38024a0c08a1dba2e90510e0a1b3c6015246304402206656347801dd620f00a86c848a06eb3369dd552b1ff905b0d4424adeb9fdb3c502201332be7725c07d84f87aefda94be83f7b3513eeeb3af7b0953e55276343a8a685a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d",
		},
		{ // public piece key
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e3802420c0899eea2e9051090b98af6024a0c0899eea2e9051090b98af6025a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d6a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dcc",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e3802420c0899eea2e9051090b98af6024a0c0899eea2e9051090b98af60252483046022100a4c53b654edb2ec19780b1c06d695c8b0bb0850edab1d2e999e784deb8f5359c0221008ad8eb26f4d6570c3dfb82faac0890b5212e30c7ad5a301220f729e3627205425a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d620b088092b8c398feffffff016a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dcc",
		},
		{ // public piece key, without piece expiration
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e38024a0c0899eea2e9051090b98af6025a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d6a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dcc",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e3802420b088092b8c398feffffff014a0c0899eea2e9051090b98af602524630440220751ae9aa91e78cf5fc858419675cb1148886acfd313c4126870d86c938675e2002206bf29b5efe3752a348446d54d3f10273bc1d582b54cbc2341db7e11508e522085a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d620b088092b8c398feffffff016a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dcc",
		},
		{ // future compatibility, "c03e01" at the end is an extra field
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e38024a0c0899eea2e9051090b98af6025a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d6a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dccc03e01",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7212200ed28abb2813e184a1e98b0f6605c4911ea468c7e8433eb583e0fca7ceac300022209566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939002a206694d2c422acd208a0072939487f6999eb9d18a44784045d87f3c67cf22746e930904e3802420b088092b8c398feffffff014a0c0899eea2e9051090b98af60252483046022100cf538a3f81f9030786bfd6d810be8b5a3c50efdfbe79ca1ce720b0a9d1359625022100e3841b64fbff97cfa8e13d9044a9b6d059b5d0914abbd1f5c792d66a5b6a50fd5a1f121d68747470733a2f2f736174656c6c6974652e6578616d706c652e636f6d620b088092b8c398feffffff016a20fd302f9f1acd1f90f5b59d8fb5d5b2d8d3d62210d4efa8647bb7a177ece96dccc03e01",
		},
	}

	for _, test := range hexes {
		unsignedBytes, err := hex.DecodeString(test.Unsigned)
		require.NoError(t, err)

		if printNewSigned {
			orderLimit := pb.OrderLimit{}
			err = pb.Unmarshal(unsignedBytes, &orderLimit)
			require.NoError(t, err)

			signed, err := signing.SignOrderLimit(ctx, signee, &orderLimit)
			require.NoError(t, err)

			signedBytes, err := pb.Marshal(signed)
			require.NoError(t, err)

			t.Log(hex.EncodeToString(signedBytes))
		}

		signedBytes, err := hex.DecodeString(test.Signed)
		require.NoError(t, err)

		orderLimit := pb.OrderLimit{}
		err = pb.Unmarshal(signedBytes, &orderLimit)
		require.NoError(t, err)

		err = signing.VerifyOrderLimitSignature(ctx, signee, &orderLimit)
		assert.NoError(t, err)

		encoded, err := signing.EncodeOrderLimit(ctx, &orderLimit)
		require.NoError(t, err)
		assert.Equal(t, unsignedBytes, encoded)
	}
}

func TestOrderVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	signer, err := identity.FullIdentityFromPEM(
		[]byte("-----BEGIN CERTIFICATE-----\nMIIBYjCCAQigAwIBAgIRAMM/5SHfNFMLl9uTAAQEoZAwCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAS/\n9wOAe42DV90jcRJMMeGe9os528RNJbMthDMkAn58KyOH87Rvlz0uCRnhhk3AbDE+\nXXHfEyed/HPFEMxJwmlGoz8wPTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYI\nKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0EAwIDSAAw\nRQIhALl9VMhM6NFnPblqOsIHOznsKr0OfQREf/+GSk/t8McsAiAxyOYg3IlB9iA0\nq/pD+qUwXuS+NFyVGOhgdNDFT3amOA==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWzCCAQGgAwIBAgIRAMfle+YJvbpRwr+FqiTrRyswCgYIKoZIzj0EAwIwEDEO\nMAwGA1UEChMFU3RvcmowIhgPMDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAw\nMFowEDEOMAwGA1UEChMFU3RvcmowWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARL\nO4n2UCp66X/MY5AzhZsfbBYOBw81Dv8V3y1BXXtbHNsUWNY8RT7r5FSTuLHsaXwq\nTwHdU05bjgnLZT/XdwqaozgwNjAOBgNVHQ8BAf8EBAMCAgQwEwYDVR0lBAwwCgYI\nKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiEA2vce\nasP0sjt6QRJNkgdV/IONJCF0IGgmsCoogCbh9ggCIA3mHgivRBId7sSAU4UUPxpB\nOOfce7bVuJlxvsnNfkkz\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBWjCCAQCgAwIBAgIQdzcArqh7Yp9aGiiJXM4+8TAKBggqhkjOPQQDAjAQMQ4w\nDAYDVQQKEwVTdG9yajAiGA8wMDAxMDEwMTAwMDAwMFoYDzAwMDEwMTAxMDAwMDAw\nWjAQMQ4wDAYDVQQKEwVTdG9yajBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM/W\nTxYhs/yGKSg8+Hb2Z/NB2KJef+fWkq7mHl7vhD9JgFwVMowMEFtKOCAhZxLBZD47\nxhYDhHBv4vrLLS+m3wGjODA2MA4GA1UdDwEB/wQEAwICBDATBgNVHSUEDDAKBggr\nBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIC+gM/sI\nXXHq5jJmolw50KKVHlqaqpdxjxJ/6x8oqTHWAiEA1w9EbqPXQ5u/oM+ODf1TBkms\nN9NfnJsY1I2A3NKEvq8=\n-----END CERTIFICATE-----\n"),
		[]byte("-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgzsFsVqt/GdqQlIIJ\nHH2VQNndv1A1fTk/35VPNzLW04ehRANCAATzXrIfcBZAHHxPdFD2PFRViRwe6eWf\nQipaF4iXQmHAW79X4mDx0BibjFfvmzurnYSlyIMZn3jp9RzbLMfnA10C\n-----END PRIVATE KEY-----\n"),
	)
	require.NoError(t, err)

	signee := signing.SignerFromFullIdentity(signer)

	type Hex struct {
		Unsigned string
		Signed   string
	}

	hexes := []Hex{
		{ // commit 385c0467
			Unsigned: "0a1068d2d6c52f5054e2d0836bf84c7174cb10e807",
			Signed:   "0a1068d2d6c52f5054e2d0836bf84c7174cb10e8071a473045022007800e9843f6ac56ae0a136406b8c685c552c7280e45761492ab521e1a27a984022100a535e3d9de1ba7778148186b319bd2857d8e2a7037a75db99b8c62eb18ed7646",
		},
	}

	for _, test := range hexes {
		unsignedBytes, err := hex.DecodeString(test.Unsigned)
		require.NoError(t, err)

		signedBytes, err := hex.DecodeString(test.Signed)
		require.NoError(t, err)

		order := pb.Order{}
		err = pb.Unmarshal(signedBytes, &order)
		require.NoError(t, err)

		err = signing.VerifyOrderSignature(ctx, signee, &order)
		assert.NoError(t, err)

		encoded, err := signing.EncodeOrder(ctx, &order)
		require.NoError(t, err)
		assert.Equal(t, unsignedBytes, encoded)
	}
}

func TestUplinkOrderVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	publicKeyBytes, _ := hex.DecodeString("01eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")
	privateKeyBytes, _ := hex.DecodeString("afefcccadb3d17b1f241b7c83f88c088b54c01b5a25409c13cbeca6bfa22b06901eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")

	publicKey, err := storj.PiecePublicKeyFromBytes(publicKeyBytes)
	require.NoError(t, err)
	privateKey, err := storj.PiecePrivateKeyFromBytes(privateKeyBytes)
	require.NoError(t, err)
	_ = privateKey

	type Hex struct {
		Unsigned string
		Signed   string
		Invalid  bool
	}

	hexes := []Hex{
		{
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7210e807",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7210e8071a4017871739c3d458737bf24bf214a7387552b18ad75afc3636974cb0d768901a85446954d59a291dde7fde0c648a242863891f543121d4633778c5b6057e62e607",
		},
		{ // future compatibility, "c03e01" at the end is an extra field
			Unsigned: "0a1052fdfc072182654f163f5f0f9a621d7210e807c03e01",
			Signed:   "0a1052fdfc072182654f163f5f0f9a621d7210e8071a40d684bccfb6494e9228cd564241183956af36af6c0ce0f49ec115bb15deaf1300f01cbbe6b3894a0b37e0c5fdd28c973d33579b0209650aa0eb80431bfd164f0dc03e01",
			Invalid:  true,
		},
	}

	for _, test := range hexes {
		unsignedBytes, err := hex.DecodeString(test.Unsigned)
		require.NoError(t, err)

		if printNewSigned {
			order := pb.Order{}
			err = pb.Unmarshal(unsignedBytes, &order)
			require.NoError(t, err)

			signed, err := signing.SignUplinkOrder(ctx, privateKey, &order)
			require.NoError(t, err)

			signedBytes, err := pb.Marshal(signed)
			require.NoError(t, err)

			t.Log(hex.EncodeToString(signedBytes))
		}

		signedBytes, err := hex.DecodeString(test.Signed)
		require.NoError(t, err)

		order := pb.Order{}
		err = pb.Unmarshal(signedBytes, &order)
		require.NoError(t, err)

		err = signing.VerifyUplinkOrderSignature(ctx, publicKey, &order)
		if test.Invalid {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		encoded, err := signing.EncodeOrder(ctx, &order)
		require.NoError(t, err)
		assert.Equal(t, unsignedBytes, encoded)
	}
}

func TestPieceHashVerification(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	publicKeyBytes, _ := hex.DecodeString("01eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")
	privateKeyBytes, _ := hex.DecodeString("afefcccadb3d17b1f241b7c83f88c088b54c01b5a25409c13cbeca6bfa22b06901eaebcb418cd629d4c01f365f33006c9de3ce70cf04da76c39cdc993f48fe53")

	publicKey, err := storj.PiecePublicKeyFromBytes(publicKeyBytes)
	require.NoError(t, err)
	privateKey, err := storj.PiecePrivateKeyFromBytes(privateKeyBytes)
	require.NoError(t, err)
	_ = privateKey

	type Hex struct {
		Unsigned string
		Signed   string
		Invalid  bool
	}

	hexes := []Hex{
		{
			Unsigned: "0a2052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649122081855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f699920e8072a0c08ba92a3e90510c89afe8202",
			Signed:   "0a2052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649122081855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f69991a40757ff5203925e02c246babdd91c9321265a158d19c99258493fe5cb6482d4bbbb97dea35227ba7b693a3c878e47d8392fc78388e225b541b98c799be7fce3c0720e8072a0c08ba92a3e90510c89afe8202",
		},
		{ // future compatibility, "c03e01" at the end is an extra field
			Unsigned: "0a2052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649122081855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f699920e8072a0c08ba92a3e90510c89afe8202c03e01",
			Signed:   "0a2052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649122081855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f69991a40e624c2fb12d7d3aa9869ba14a6f0c9fe0edefa046eada2126d61d6d3915515fcfe47763fb9a29575997949d2b824e079bc91d88d56a89b963381cf176d00500e20e8072a0c08ba92a3e90510c89afe8202c03e01",
			Invalid:  true,
		},
	}

	for _, test := range hexes {
		unsignedBytes, err := hex.DecodeString(test.Unsigned)
		require.NoError(t, err)

		if printNewSigned {
			hash := pb.PieceHash{}
			err = pb.Unmarshal(unsignedBytes, &hash)
			require.NoError(t, err)

			signed, err := signing.SignUplinkPieceHash(ctx, privateKey, &hash)
			require.NoError(t, err)

			signedBytes, err := pb.Marshal(signed)
			require.NoError(t, err)

			t.Log(hex.EncodeToString(signedBytes))
		}

		signedBytes, err := hex.DecodeString(test.Signed)
		require.NoError(t, err)

		hash := pb.PieceHash{}
		err = pb.Unmarshal(signedBytes, &hash)
		require.NoError(t, err)

		err = signing.VerifyUplinkPieceHashSignature(ctx, publicKey, &hash)
		if test.Invalid {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		encoded, err := signing.EncodePieceHash(ctx, &hash)
		require.NoError(t, err)
		assert.Equal(t, unsignedBytes, encoded)
	}
}

func TestSignExitCompleted(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	satIdentity, err := testidentity.NewTestIdentity(ctx)
	nodeID := testrand.NodeID()
	require.NoError(t, err)

	finishedAt := time.Now().UTC()
	signer := signing.SignerFromFullIdentity(satIdentity)
	signee := signing.SigneeFromPeerIdentity(satIdentity.PeerIdentity())

	unsigned := &pb.ExitCompleted{
		SatelliteId: satIdentity.ID,
		NodeId:      nodeID,
		Completed:   finishedAt,
	}
	signed, err := signing.SignExitCompleted(ctx, signer, unsigned)
	require.NoError(t, err)

	err = signing.VerifyExitCompleted(ctx, signee, signed)
	require.NoError(t, err)

	signed.SatelliteId = testrand.NodeID()

	err = signing.VerifyExitCompleted(ctx, signee, signed)
	require.Error(t, err)
}

func TestSignExitFailed(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	satIdentity, err := testidentity.NewTestIdentity(ctx)
	nodeID := testrand.NodeID()
	require.NoError(t, err)

	finishedAt := time.Now().UTC()
	signer := signing.SignerFromFullIdentity(satIdentity)
	signee := signing.SigneeFromPeerIdentity(satIdentity.PeerIdentity())

	unsigned := &pb.ExitFailed{
		SatelliteId: satIdentity.ID,
		NodeId:      nodeID,
		Failed:      finishedAt,
		Reason:      pb.ExitFailed_INACTIVE_TIMEFRAME_EXCEEDED,
	}
	signed, err := signing.SignExitFailed(ctx, signer, unsigned)
	require.NoError(t, err)

	err = signing.VerifyExitFailed(ctx, signee, signed)
	require.NoError(t, err)

	signed.Reason = pb.ExitFailed_OVERALL_FAILURE_PERCENTAGE_EXCEEDED

	err = signing.VerifyExitFailed(ctx, signee, signed)
	require.Error(t, err)
}
