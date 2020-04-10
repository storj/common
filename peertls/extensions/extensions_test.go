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
		revocation extensions.Revocation
	}{
		{revocation: extensions.Revocation{}},
		{revocation: extensions.Revocation{Timestamp: 1}},
		{revocation: extensions.Revocation{Timestamp: 9223372036854775807}},
		{revocation: extensions.Revocation{KeyHash: []byte{1, 2, 3}}},
		{revocation: extensions.Revocation{Signature: []byte{5, 4, 3}}},
		{revocation: extensions.Revocation{
			Timestamp: 9223372036854775807,
			KeyHash:   []byte{5, 4, 3}},
		},
		{revocation: extensions.Revocation{
			Timestamp: 9223372036854775807,
			Signature: []byte{5, 4, 3}},
		},
		{revocation: extensions.Revocation{
			KeyHash:   []byte{1, 2, 3},
			Signature: []byte{5, 4, 3}},
		},
		{revocation: extensions.Revocation{
			Timestamp: testrand.Int63n(9223372036854775807),
			KeyHash:   testrand.BytesInt(testrand.Intn(600000)),
			Signature: testrand.BytesInt(testrand.Intn(500000))},
		},
	} {
		gobEncoded := new(bytes.Buffer)
		encoder := gob.NewEncoder(gobEncoded)
		err := encoder.Encode(tt.revocation)
		require.NoError(t, err)

		customEncoded, err := tt.revocation.Marshal()
		require.NoError(t, err)

		// compare gob marshaler output with our marshaler
		require.Equal(t, gobEncoded.Bytes(), customEncoded)

		revocatationDocodeGob := &extensions.Revocation{}
		gobDecoder := gob.NewDecoder(bytes.NewBuffer(customEncoded))
		err = gobDecoder.Decode(revocatationDocodeGob)
		require.NoError(t, err)
		require.Equal(t, tt.revocation, *revocatationDocodeGob)

		// unmarshal data from gob marshaler with our marshaler
		revocatationGob := extensions.Revocation{}
		err = revocatationGob.Unmarshal(gobEncoded.Bytes())
		require.NoError(t, err)
		require.Equal(t, tt.revocation, revocatationGob)
		require.Equal(t, *revocatationDocodeGob, revocatationGob)

		// unmarshal data from our marshaler with our marshaler
		revocatationCustom := extensions.Revocation{}
		err = revocatationCustom.Unmarshal(customEncoded)
		require.NoError(t, err)
		require.Equal(t, tt.revocation, revocatationCustom)
		require.Equal(t, *revocatationDocodeGob, revocatationCustom)
	}
}

func TestRevocationMarshalingInvalid(t *testing.T) {
	gobEncoded := new(bytes.Buffer)
	encoder := gob.NewEncoder(gobEncoded)
	// encode different object
	err := encoder.Encode(extensions.Options{})
	require.NoError(t, err)

	revocatationCustom := extensions.Revocation{}
	err = revocatationCustom.Unmarshal(gobEncoded.Bytes())
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
