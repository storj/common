// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package extensions_test

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/identity"
	"storj.io/common/peertls/extensions"
	"storj.io/common/peertls/testpeertls"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestHandlers_Register(t *testing.T) {
	var (
		handlers = extensions.HandlerFactories{}
		ids      []*extensions.ExtensionID
		opts     []*extensions.Options
		exts     []pkix.Extension
		chains   [][][]*x509.Certificate
	)

	for idx := 0; idx < 5; idx++ {
		i := idx

		ids = append(ids, &extensions.ExtensionID{2, 999, 999, i})
		opts = append(opts, &extensions.Options{})
		exts = append(exts, pkix.Extension{Id: *ids[i]})

		_, chain, err := testpeertls.NewCertChain(2, storj.LatestIDVersion().Number)
		require.NoError(t, err)
		chains = append(chains, identity.ToChains(chain))

		testHandler := extensions.NewHandlerFactory(
			ids[i],
			func(opt *extensions.Options) extensions.HandlerFunc {
				assert.Equal(t, opts[i], opt)
				assert.NotNil(t, opt)

				return func(ext pkix.Extension, chain [][]*x509.Certificate) error {
					assert.NotNil(t, ext)
					assert.Equal(t, exts[i], ext)

					assert.NotNil(t, ext.Id)
					assert.Equal(t, *ids[i], ext.Id)

					assert.NotNil(t, chain)
					assert.Equal(t, chains[i], chain)
					return errs.New(strconv.Itoa(i))
				}
			},
		)
		handlers.Register(testHandler)

		err = handlers[i].NewHandlerFunc(opts[i])(exts[i], chains[i])
		assert.Errorf(t, err, strconv.Itoa(i))
	}
}

func TestHandlers_WithOptions(t *testing.T) {
	var (
		handlers = extensions.HandlerFactories{}
		ids      []*extensions.ExtensionID
		opts     []*extensions.Options
		exts     []pkix.Extension
		chains   [][][]*x509.Certificate
	)

	for idx := 0; idx < 5; idx++ {
		i := idx

		ids = append(ids, &extensions.ExtensionID{2, 999, 999, i})
		opts = append(opts, &extensions.Options{})
		exts = append(exts, pkix.Extension{Id: *ids[i]})

		_, chain, err := testpeertls.NewCertChain(2, storj.LatestIDVersion().Number)
		require.NoError(t, err)
		chains = append(chains, identity.ToChains(chain))

		testHandler := extensions.NewHandlerFactory(
			ids[i],
			func(opt *extensions.Options) extensions.HandlerFunc {
				assert.Equal(t, opts[i], opt)
				assert.NotNil(t, opt)

				return func(ext pkix.Extension, chain [][]*x509.Certificate) error {
					assert.NotNil(t, ext)
					assert.Equal(t, exts[i], ext)

					assert.NotNil(t, ext.Id)
					assert.Equal(t, *ids[i], ext.Id)

					assert.NotNil(t, chain)
					assert.Equal(t, chains[i], chain)
					return errs.New(strconv.Itoa(i))
				}
			},
		)
		handlers.Register(testHandler)

		handlerFuncMap := handlers.WithOptions(&extensions.Options{})

		id := handlers[i].ID()
		require.NotNil(t, id)

		handleFunc, ok := handlerFuncMap[id]
		assert.True(t, ok)
		assert.NotNil(t, handleFunc)
	}
}

func TestRevocationMarshaling(t *testing.T) {
	for _, tt := range []struct {
		// gob is the older version of Gob encoding
		gobbytes   []byte
		revocation extensions.Revocation
	}{
		{
			gobbytes:   []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x3, 0xff, 0x82, 0x0},
			revocation: extensions.Revocation{},
		}, {
			gobbytes:   []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x5, 0xff, 0x82, 0x1, 0x2, 0x0},
			revocation: extensions.Revocation{Timestamp: 1},
		}, {
			gobbytes:   []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0xd, 0xff, 0x82, 0x1, 0xf8, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0x0},
			revocation: extensions.Revocation{Timestamp: 9223372036854775807},
		}, {
			gobbytes:   []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x8, 0xff, 0x82, 0x2, 0x3, 0x1, 0x2, 0x3, 0x0},
			revocation: extensions.Revocation{KeyHash: []byte{1, 2, 3}},
		}, {
			gobbytes:   []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x8, 0xff, 0x82, 0x3, 0x3, 0x5, 0x4, 0x3, 0x0},
			revocation: extensions.Revocation{Signature: []byte{5, 4, 3}},
		}, {
			gobbytes: []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x12, 0xff, 0x82, 0x1, 0xf8, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0x1, 0x3, 0x5, 0x4, 0x3, 0x0},
			revocation: extensions.Revocation{
				Timestamp: 9223372036854775807,
				KeyHash:   []byte{5, 4, 3}},
		}, {
			gobbytes: []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0x12, 0xff, 0x82, 0x1, 0xf8, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0x2, 0x3, 0x5, 0x4, 0x3, 0x0},
			revocation: extensions.Revocation{
				Timestamp: 9223372036854775807,
				Signature: []byte{5, 4, 3}},
		}, {
			gobbytes: []byte{0x40, 0xff, 0x81, 0x3, 0x1, 0x1, 0xa, 0x52, 0x65, 0x76, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1, 0xff, 0x82, 0x0, 0x1, 0x3, 0x1, 0x9, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1, 0x4, 0x0, 0x1, 0x7, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x1, 0xa, 0x0, 0x1, 0x9, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1, 0xa, 0x0, 0x0, 0x0, 0xd, 0xff, 0x82, 0x2, 0x3, 0x1, 0x2, 0x3, 0x1, 0x3, 0x5, 0x4, 0x3, 0x0},
			revocation: extensions.Revocation{
				KeyHash:   []byte{1, 2, 3},
				Signature: []byte{5, 4, 3}},
		}, {
			gobbytes: nil, // skip the encoding test for this
			revocation: extensions.Revocation{
				Timestamp: testrand.Int63n(9223372036854775807),
				KeyHash:   testrand.BytesInt(testrand.Intn(600000)),
				Signature: testrand.BytesInt(testrand.Intn(500000))},
		},
	} {
		customEncoded, err := tt.revocation.Marshal()
		require.NoError(t, err)

		if tt.gobbytes != nil {
			// compare gob marshaler output with our marshaler
			require.Equal(t, tt.gobbytes, customEncoded)
		}

		revocatationDocodeGob := &extensions.Revocation{}
		gobDecoder := gob.NewDecoder(bytes.NewBuffer(customEncoded))
		err = gobDecoder.Decode(revocatationDocodeGob)
		require.NoError(t, err)
		require.Equal(t, tt.revocation, *revocatationDocodeGob)

		// unmarshal data from our marshaler with our marshaler
		revocatationCustom := extensions.Revocation{}
		err = revocatationCustom.Unmarshal(customEncoded)
		require.NoError(t, err)
		require.Equal(t, tt.revocation, revocatationCustom)
		require.Equal(t, *revocatationDocodeGob, revocatationCustom)
	}
}

func TestRevocationMarshalingInvalid(t *testing.T) {
	// the gob bytes were generated using: Go 1.21
	// var b bytes.Buffer
	// encoder := gob.NewEncoder(&b)
	// encoder.Encode(extensions.Options{})
	// encoded = b.Bytes()

	var encoded = []byte("M\x7f\x03\x01\x01\aOptions\x01\xff\x80\x00\x01\x03\x01\x0fPeerCAWhitelist\x01\xff\xa4\x00\x01\fRevocationDB\x01\x10\x00\x01\x0ePeerIDVersions\x01\f\x00\x00\x00\"\xff\xa3\x02\x01\x01\x13[]*x509.Certificate\x01\xff\xa4\x00\x01\xff\x82\x00\x00\xfe\x03h\xff\x81\x03\x01\x02\xff\x82\x00\x01,\x01\x03Raw\x01\n\x00\x01\x11RawTBSCertificate\x01\n\x00\x01\x17RawSubjectPublicKeyInfo\x01\n\x00\x01\nRawSubject\x01\n\x00\x01\tRawIssuer\x01\n\x00\x01\tSignature\x01\n\x00\x01\x12SignatureAlgorithm\x01\x04\x00\x01\x12PublicKeyAlgorithm\x01\x04\x00\x01\tPublicKey\x01\x10\x00\x01\aVersion\x01\x04\x00\x01\fSerialNumber\x01\xff\x84\x00\x01\x06Issuer\x01\xff\x86\x00\x01\aSubject\x01\xff\x86\x00\x01\tNotBefore\x01\xff\x90\x00\x01\bNotAfter\x01\xff\x90\x00\x01\bKeyUsage\x01\x04\x00\x01\nExtensions\x01\xff\x94\x00\x01\x0fExtraExtensions\x01\xff\x94\x00\x01\x1bUnhandledCriticalExtensions\x01\xff\x96\x00\x01\vExtKeyUsage\x01\xff\x98\x00\x01\x12UnknownExtKeyUsage\x01\xff\x96\x00\x01\x15BasicConstraintsValid\x01\x02\x00\x01\x04IsCA\x01\x02\x00\x01\nMaxPathLen\x01\x04\x00\x01\x0eMaxPathLenZero\x01\x02\x00\x01\fSubjectKeyId\x01\n\x00\x01\x0eAuthorityKeyId\x01\n\x00\x01\nOCSPServer\x01\xff\x88\x00\x01\x15IssuingCertificateURL\x01\xff\x88\x00\x01\bDNSNames\x01\xff\x88\x00\x01\x0eEmailAddresses\x01\xff\x88\x00\x01\vIPAddresses\x01\xff\x9a\x00\x01\x04URIs\x01\xff\x9e\x00\x01\x1bPermittedDNSDomainsCritical\x01\x02\x00\x01\x13PermittedDNSDomains\x01\xff\x88\x00\x01\x12ExcludedDNSDomains\x01\xff\x88\x00\x01\x11PermittedIPRanges\x01\xff\xa2\x00\x01\x10ExcludedIPRanges\x01\xff\xa2\x00\x01\x17PermittedEmailAddresses\x01\xff\x88\x00\x01\x16ExcludedEmailAddresses\x01\xff\x88\x00\x01\x13PermittedURIDomains\x01\xff\x88\x00\x01\x12ExcludedURIDomains\x01\xff\x88\x00\x01\x15CRLDistributionPoints\x01\xff\x88\x00\x01\x11PolicyIdentifiers\x01\xff\x96\x00\x00\x00\n\xff\x83\x05\x01\x02\xff\xa6\x00\x00\x00\xff\xc3\xff\x85\x03\x01\x01\x04Name\x01\xff\x86\x00\x01\v\x01\aCountry\x01\xff\x88\x00\x01\fOrganization\x01\xff\x88\x00\x01\x12OrganizationalUnit\x01\xff\x88\x00\x01\bLocality\x01\xff\x88\x00\x01\bProvince\x01\xff\x88\x00\x01\rStreetAddress\x01\xff\x88\x00\x01\nPostalCode\x01\xff\x88\x00\x01\fSerialNumber\x01\f\x00\x01\nCommonName\x01\f\x00\x01\x05Names\x01\xff\x8e\x00\x01\nExtraNames\x01\xff\x8e\x00\x00\x00\x16\xff\x87\x02\x01\x01\b[]string\x01\xff\x88\x00\x01\f\x00\x00+\xff\x8d\x02\x01\x01\x1c[]pkix.AttributeTypeAndValue\x01\xff\x8e\x00\x01\xff\x8a\x00\x007\xff\x89\x03\x01\x01\x15AttributeTypeAndValue\x01\xff\x8a\x00\x01\x02\x01\x04Type\x01\xff\x8c\x00\x01\x05Value\x01\x10\x00\x00\x00\x1e\xff\x8b\x02\x01\x01\x10ObjectIdentifier\x01\xff\x8c\x00\x01\x04\x00\x00\x10\xff\x8f\x05\x01\x01\x04Time\x01\xff\x90\x00\x00\x00\x1f\xff\x93\x02\x01\x01\x10[]pkix.Extension\x01\xff\x94\x00\x01\xff\x92\x00\x006\xff\x91\x03\x01\x01\tExtension\x01\xff\x92\x00\x01\x03\x01\x02Id\x01\xff\x8c\x00\x01\bCritical\x01\x02\x00\x01\x05Value\x01\n\x00\x00\x00&\xff\x95\x02\x01\x01\x17[]asn1.ObjectIdentifier\x01\xff\x96\x00\x01\xff\x8c\x00\x00 \xff\x97\x02\x01\x01\x12[]x509.ExtKeyUsage\x01\xff\x98\x00\x01\x04\x00\x00\x16\xff\x99\x02\x01\x01\b[]net.IP\x01\xff\x9a\x00\x01\n\x00\x00\x19\xff\x9d\x02\x01\x01\n[]*url.URL\x01\xff\x9e\x00\x01\xff\x9c\x00\x00\n\xff\x9b\x06\x01\x02\xff\xa8\x00\x00\x00\x14\xff\xa9\x03\x01\x01\bUserinfo\x01\xff\xaa\x00\x00\x00\x1b\xff\xa1\x02\x01\x01\f[]*net.IPNet\x01\xff\xa2\x00\x01\xff\xa0\x00\x00\x1c\xff\x9f\x03\x01\x02\xff\xa0\x00\x01\x02\x01\x02IP\x01\n\x00\x01\x04Mask\x01\n\x00\x00\x00\x03\xff\x80\x00")
	revocatationCustom := extensions.Revocation{}
	err := revocatationCustom.Unmarshal(encoded)
	require.Error(t, err)

	// try to unmarshal random bytes
	revocatationCustom = extensions.Revocation{}
	err = revocatationCustom.Unmarshal(testrand.BytesInt(10000))
	require.Error(t, err)
}

func TestRevocationDecoderCrashers(t *testing.T) {
	crashers := []string{
		"@\xff\x81\x03\x01\x01\nRevocation\x01\xff\x82\x00\x01\x03\x01\tTimestamp\x01\x04\x00\x01\aKeyHash\x01\n\x00\x01\tSignature\x01\n\x00\x00\x00\r\xff\x82\x02\xf8000000000",
	}

	for _, crasher := range crashers {
		rev := extensions.Revocation{}
		_ = rev.Unmarshal([]byte(crasher))
	}
}
