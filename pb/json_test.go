// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package pb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	t.Run("CipherSuite", func(t *testing.T) {
		cs := CipherSuite_ENC_AESGCM
		bs, err := json.Marshal(cs)
		require.NoError(t, err)
		require.Equal(t, string(bs), `"ENC_AESGCM"`)
	})

	t.Run("EncryptionAccess_StoreEntry", func(t *testing.T) {
		bs, err := json.Marshal(&EncryptionAccess_StoreEntry{
			Bucket:          []byte("bucket"),
			UnencryptedPath: []byte("unenc"),
			EncryptedPath:   []byte("enc"),
			Key:             []byte("key"),
			PathCipher:      CipherSuite_ENC_AESGCM,
			MetadataCipher:  CipherSuite_ENC_AESGCM,
		})
		require.NoError(t, err)
		require.Equal(t, string(bs), `{"bucket":"bucket","unencrypted_path":"unenc","encrypted_path":"ZW5j","key":"a2V5","path_cipher":"ENC_AESGCM","metadata_cipher":"ENC_AESGCM"}`)
	})
}
