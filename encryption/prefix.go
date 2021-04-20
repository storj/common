// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"github.com/zeebo/errs"

	"storj.io/common/paths"
	"storj.io/common/storj"
)

// PrefixInfo is a helper type that contains all of the encrypted and unencrypted paths related
// to some path and its parent. It includes the cipher that was used to encrypt and decrypt
// the paths and what bucket it is in.
type PrefixInfo struct {
	Bucket string
	Cipher storj.CipherSuite

	PathUnenc paths.Unencrypted
	PathEnc   paths.Encrypted
	PathKey   storj.Key

	ParentUnenc paths.Unencrypted
	ParentEnc   paths.Encrypted
	ParentKey   storj.Key
}

// GetPrefixInfo returns the PrefixInfo for some unencrypted path inside of a bucket.
func GetPrefixInfo(bucket string, path paths.Unencrypted, store *Store) (pi *PrefixInfo, err error) {
	_, consumed, base := store.LookupUnencrypted(bucket, path)
	if base == nil {
		return nil, ErrMissingDecryptionBase.New("%q/%q", bucket, path)
	}

	remaining, ok := path.Consume(consumed)
	if !ok {
		return nil, errs.New("unable to encrypt bucket path: %q/%q", bucket, path)
	}

	// if we're using the default base (meaning the default key), we need
	// to include the bucket name in the path derivation.
	key := &base.Key
	if base.Default {
		key, err = derivePathKeyComponent(key, bucket)
		if err != nil {
			return nil, errs.Wrap(err)
		}
	}

	pi = &PrefixInfo{
		Bucket: bucket,
		Cipher: base.PathCipher,

		PathUnenc: path,
		PathKey:   *key,

		ParentKey: *key,
	}

	var componentUnenc string
	var componentEnc string
	var pathEnc pathBuilder
	var parentEnc pathBuilder
	var parentUnenc pathBuilder

	for iter, i := remaining.Iterator(), 0; !iter.Done(); i++ {
		if i > 0 {
			parentEnc.append(componentEnc)
			parentUnenc.append(componentUnenc)
			pi.ParentKey = *key
		}

		componentUnenc = iter.Next()

		componentEnc, err = encryptPathComponent(componentUnenc, base.PathCipher, key)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		key, err = derivePathKeyComponent(key, componentUnenc)
		if err != nil {
			return nil, errs.Wrap(err)
		}

		pathEnc.append(componentEnc)
		pi.PathKey = *key
	}

	pi.PathEnc = paths.NewEncrypted(pathEnc.String())
	pi.ParentUnenc = paths.NewUnencrypted(parentUnenc.String())
	pi.ParentEnc = paths.NewEncrypted(parentEnc.String())

	return pi, nil
}
