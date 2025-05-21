// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/paths"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func newStore(key storj.Key, pathCipher storj.CipherSuite) *Store {
	store := NewStore()
	if err := store.AddWithCipher("bucket", paths.Unencrypted{}, paths.Encrypted{}, key, pathCipher); err != nil {
		panic(err)
	}
	return store
}

func TestStoreEncryption(t *testing.T) {
	forAllCiphers(func(cipher storj.CipherSuite) {
		for i, rawPath := range []string{
			"",
			"/",
			"//",
			"file.txt",
			"file.txt/",
			"fold1/file.txt",
			"fold1/fold2/file.txt",
			"/fold1/fold2/fold3/file.txt",
			"/fold1/fold2/fold3/file.txt/",
		} {
			errTag := fmt.Sprintf("test:%d path:%q cipher:%v", i, rawPath, cipher)

			store := newStore(testrand.Key(), cipher)
			path := paths.NewUnencrypted(rawPath)

			encPath, err := EncryptPathWithStoreCipher("bucket", path, store)
			if !assert.NoError(t, err, errTag) {
				continue
			}
			if cipher != storj.EncNull {
				if !assert.True(t, !strings.HasSuffix(encPath.Raw(), "/"), errTag) {
					continue
				}
			}

			decPath, err := DecryptPathWithStoreCipher("bucket", encPath, store)
			if !assert.NoError(t, err, errTag) {
				continue
			}

			assert.Equal(t, rawPath, decPath.Raw(), errTag)
		}
	})
}

func TestStorePrefixEncryption(t *testing.T) {
	forAllCiphers(func(cipher storj.CipherSuite) {
		for i, rawPath := range []string{
			"",
			"/",
			"//",
			"file.txt",
			"file.txt/",
			"fold1/file.txt",
			"fold1/fold2/file.txt",
			"/fold1/fold2/fold3/file.txt",
			"/fold1/fold2/fold3/file.txt/",
		} {
			errTag := fmt.Sprintf("test:%d path:%q cipher:%v", i, rawPath, cipher)

			store := newStore(testrand.Key(), cipher)
			path := paths.NewUnencrypted(rawPath)

			encPath, err := EncryptPrefixWithStoreCipher("bucket", path, store)
			if !assert.NoError(t, err, errTag) {
				continue
			}
			if !assert.Equal(t,
				strings.HasSuffix(path.Raw(), "/"),
				strings.HasSuffix(encPath.Raw(), "/"),
				errTag) {
				continue
			}

			decPath, err := DecryptPathWithStoreCipher("bucket", encPath, store)
			if !assert.NoError(t, err, errTag) {
				continue
			}

			assert.Equal(t, rawPath, decPath.Raw(), errTag)
		}
	})
}

func TestStoreEncryption_BucketRoot(t *testing.T) {
	forAllCiphers(func(cipher storj.CipherSuite) {
		for i, rawPath := range []string{
			"",
			"/",
			"//",
			"file.txt",
			"file.txt/",
			"fold1/file.txt",
			"fold1/fold2/file.txt",
			"/fold1/fold2/fold3/file.txt",
			"/fold1/fold2/fold3/file.txt/",
		} {
			errTag := fmt.Sprintf("test:%d path:%q cipher:%v", i, rawPath, cipher)

			dk := testrand.Key()
			rootStore := NewStore()
			rootStore.SetDefaultKey(&dk)
			rootStore.SetDefaultPathCipher(cipher)

			bucketStore := NewStore()
			bucketKey, err := DerivePathKey("bucket", paths.Unencrypted{}, rootStore)
			if !assert.NoError(t, err, errTag) {
				continue
			}
			err = bucketStore.AddWithCipher("bucket", paths.Unencrypted{}, paths.Encrypted{}, *bucketKey, cipher)
			if !assert.NoError(t, err, errTag) {
				continue
			}

			path := paths.NewUnencrypted(rawPath)

			rootEncPath, err := EncryptPathWithStoreCipher("bucket", path, rootStore)
			if !assert.NoError(t, err, errTag) {
				continue
			}

			bucketEncPath, err := EncryptPathWithStoreCipher("bucket", path, bucketStore)
			if !assert.NoError(t, err, errTag) {
				continue
			}

			assert.Equal(t, rootEncPath, bucketEncPath, errTag)
		}
	})
}

func TestStoreEncryption_MultipleBases(t *testing.T) {
	forAllCiphers(func(cipher storj.CipherSuite) {
		for _, rawPath := range []string{
			"",
			"/",
			"//",
			"file.txt",
			"file.txt/",
			"fold1/file.txt",
			"fold1/fold2/file.txt",
			"/fold1/fold2/fold3/file.txt",
			"/fold1/fold2/fold3/file.txt/",
		} {

			var pb pathBuilder
			for iter := paths.NewIterator(rawPath); !iter.Done(); {
				pb.append(iter.Next())
				prefix := pb.Unencrypted()
				rawPath := rawPath

				t.Run(fmt.Sprintf("path:%q prefix:%q cipher:%v", rawPath, prefix, cipher), func(t *testing.T) {
					dk := testrand.Key()

					rootStore := NewStore()
					rootStore.SetDefaultKey(&dk)
					rootStore.SetDefaultPathCipher(cipher)

					prefixStore := NewStore()
					prefixStore.SetDefaultKey(&dk)
					prefixStore.SetDefaultPathCipher(cipher)

					prefixKey, err := DerivePathKey("bucket", prefix, rootStore)
					require.NoError(t, err)

					encPrefix, err := EncryptPath("bucket", prefix, cipher, rootStore)
					require.NoError(t, err)

					err = prefixStore.AddWithCipher("bucket", prefix, encPrefix, *prefixKey, cipher)
					require.NoError(t, err)

					path := paths.NewUnencrypted(rawPath)

					rootEncPath, err := EncryptPathWithStoreCipher("bucket", path, rootStore)
					require.NoError(t, err)

					prefixEncPath, err := EncryptPathWithStoreCipher("bucket", path, prefixStore)
					require.NoError(t, err)

					require.Equal(t, rootEncPath, prefixEncPath)
				})
			}
		}
	})
}

func forAllCiphers(test func(cipher storj.CipherSuite)) {
	for _, cipher := range []storj.CipherSuite{
		storj.EncNull,
		storj.EncAESGCM,
		storj.EncSecretBox,
	} {
		test(cipher)
	}
}

func TestSegmentEncoding(t *testing.T) {
	segments := [][]byte{
		{},
		{'a'},
		{0},
		{'/'},
		{'a', 'b', 'c', 'd', '1', '2', '3', '4', '5'},
		{'/', '/', '/', '/', '/'},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{'a', '/', 'a', '2', 'a', 'a', 0, '1', 'b', 255},
		{'/', '/', 'a', 0, 'a', 'a', 0, '1', 'b', 'g', 'a', 'b', '/'},
		{0, '/', 'a', '0', 'a', 'a', 0, '1', 'b', 'g', 'a', 'b', 0},
	}

	// additional random segments
	for range 20 {
		segments = append(segments, testrand.BytesInt(testrand.Intn(256)))
	}

	for i, segment := range segments {
		encoded := encodeSegment(segment)
		require.Equal(t, -1, bytes.IndexByte(encoded, 0))
		require.Equal(t, -1, bytes.IndexByte(encoded, 255))
		require.Equal(t, -1, bytes.IndexByte(encoded, '/'))

		decoded, err := decodeSegment(encoded)
		require.NoError(t, err, "#%d", i)
		require.Equal(t, segment, decoded, "#%d", i)
	}
}

func TestInvalidSegmentDecoding(t *testing.T) {
	encoded := []byte{3, 4, 5, 6, 7}
	// first byte should be '\x01' or '\x02'
	_, err := decodeSegment(encoded)
	require.Error(t, err)
	assert.True(t, ErrDecryptFailed.Has(err), "invalid error class")
}

func TestValidateEncodedSegment(t *testing.T) {
	// all segments should be invalid
	encodedSegments := [][]byte{
		{},
		{1, 1},
		{2},
		{2, 0},
		{2, '\xff'},
		{2, '\x2f'},
		{2, escapeSlash, '3'},
		{2, escapeFF, '3'},
		{2, escape01, '3'},
		{3, 4, 4, 4},
	}

	for i, segment := range encodedSegments {
		_, err := decodeSegment(segment)
		require.Error(t, err, "#%d", i)
		assert.True(t, ErrDecryptFailed.Has(err), "invalid error class #%d", i)
	}
}

func TestEncodingDecodingStress(t *testing.T) {
	specials := [...]byte{0x00, 0x01, 0x02, 0x03, 'A', '/', '\\', 0x2d, 0x2e, 0x2f, 0xfd, 0xfe, 0xff}
	const n = len(specials)
	for i := range n * n * n {
		var segment [3]byte
		segment[0] = specials[i%n]
		segment[1] = specials[i/n%n]
		segment[2] = specials[i/(n*n)%n]

		_ = encodeSegment(segment[:])
		_, _ = decodeSegment(segment[:])
	}

	// random segments
	for range 20 {
		segment := testrand.BytesInt(testrand.Intn(256))
		_ = encodeSegment(segment)
		_, _ = decodeSegment(segment)
	}
}

func TestDecryptPath_EncryptionBypass(t *testing.T) {
	encStore := NewStore()
	encStore.SetDefaultKey(&storj.Key{})
	encStore.SetDefaultPathCipher(storj.EncAESGCM)

	bucketName := "test-bucket"

	filePaths := []string{
		"a", "aa", "b", "bb", "c",
		"a/xa", "a/xaa", "a/xb", "a/xbb", "a/xc",
		"b/ya", "b/yaa", "b/yb", "b/ybb", "b/yc",
	}

	for _, path := range filePaths {
		encStore.EncryptionBypass = false
		encryptedPath, err := EncryptPathWithStoreCipher(bucketName, paths.NewUnencrypted(path), encStore)
		require.NoError(t, err)

		var expectedPath, next string
		iterator := encryptedPath.Iterator()
		for !iterator.Done() {
			next = iterator.Next()
			expectedPath += base64.URLEncoding.EncodeToString([]byte(next)) + "/"
		}
		expectedPath = strings.TrimRight(expectedPath, "/")

		encStore.EncryptionBypass = true
		actualPath, err := DecryptPathWithStoreCipher(bucketName, encryptedPath, encStore)
		require.NoError(t, err)

		require.Equal(t, paths.NewUnencrypted(expectedPath), actualPath)
	}
}

func TestEncryptPath_EncryptionBypass(t *testing.T) {
	encStore := NewStore()
	encStore.SetDefaultKey(&storj.Key{})
	encStore.SetDefaultPathCipher(storj.EncAESGCM)

	bucketName := "test-bucket"

	filePaths := []string{
		"a", "aa", "b", "bb", "c",
		"a/xa", "a/xaa", "a/xb", "a/xbb", "a/xc",
		"b/ya", "b/yaa", "b/yb", "b/ybb", "b/yc",
	}

	for _, path := range filePaths {
		encStore.EncryptionBypass = false
		encryptedPath, err := EncryptPathWithStoreCipher(bucketName, paths.NewUnencrypted(path), encStore)
		require.NoError(t, err)

		var encodedPath, next string
		iterator := encryptedPath.Iterator()
		for !iterator.Done() {
			next = iterator.Next()
			encodedPath += base64.URLEncoding.EncodeToString([]byte(next)) + "/"
		}
		encodedPath = strings.TrimRight(encodedPath, "/")

		encStore.EncryptionBypass = true
		actualPath, err := EncryptPathWithStoreCipher(bucketName, paths.NewUnencrypted(encodedPath), encStore)
		require.NoError(t, err)

		require.Equal(t, encryptedPath.String(), actualPath.String())
	}
}

func BenchmarkSegmentEncoding(b *testing.B) {
	segments := [][]byte{
		{},
		{'a'},
		{0},
		{'/'},
		{'a', 'b', 'c', 'd', '1', '2', '3', '4', '5'},

		{'/', '/', '/', '/', '/'},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},

		{'a', '/', 'a', '2', 'a', 'a', 0, '1', 'b', 255},
		{'/', '/', 'a', 0, 'a', 'a', 0, '1', 'b', 'g', 'a', 'b', '/'},
		{0, '/', 'a', '0', 'a', 'a', 0, '1', 'b', 'g', 'a', 'b', 0},
	}

	// additional random segment
	segments = append(segments, testrand.BytesInt(255))

	b.Run("Loop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, segment := range segments {
				encoded := encodeSegment(segment)
				_, _ = decodeSegment(encoded)
			}
		}
	})
	b.Run("Base64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, segment := range segments {
				encoded := base64.RawURLEncoding.EncodeToString(segment)
				_, _ = base64.RawURLEncoding.DecodeString(encoded)
			}
		}
	})
}
