// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package grant

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/macaroon"
	"storj.io/common/paths"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestRestrict(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	secret, err := macaroon.NewSecret()
	require.NoError(t, err)

	apiKey, err := macaroon.NewAPIKey(secret)
	require.NoError(t, err)

	defaultKey := testrand.Key()
	encAccess := NewEncryptionAccessWithDefaultKey(&defaultKey)
	encAccess.SetDefaultPathCipher(storj.EncNull)

	access := Access{
		APIKey:    apiKey,
		EncAccess: encAccess,
	}

	fullPermission := Permission{
		AllowDownload:                         true,
		AllowUpload:                           true,
		AllowList:                             true,
		AllowDelete:                           true,
		AllowPutObjectRetention:               true,
		AllowGetObjectRetention:               true,
		AllowPutObjectLegalHold:               true,
		AllowGetObjectLegalHold:               true,
		AllowBypassGovernanceRetention:        true,
		AllowPutBucketObjectLockConfiguration: true,
		AllowGetBucketObjectLockConfiguration: true,
	}

	action1 := macaroon.Action{
		Op:            macaroon.ActionRead,
		Time:          now,
		Bucket:        []byte("bucket"),
		EncryptedPath: []byte("prefix1/path1"),
	}

	action2 := macaroon.Action{
		Op:            macaroon.ActionRead,
		Time:          now,
		Bucket:        []byte("bucket"),
		EncryptedPath: []byte("prefix2/path2"),
	}

	// Restrict the access only with NotAfter time
	permission := fullPermission
	permission.NotAfter = time.Now().Add(2 * time.Hour)

	restricted, err := access.Restrict(permission)
	require.NoError(t, err)

	// Check that all actions are allowed and the encAccess has only the default key
	assert.NoError(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action1, nil))
	assert.NoError(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action2, nil))
	assert.Equal(t, &defaultKey, restricted.EncAccess.Store.GetDefaultKey())

	_, _, base := restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix1/path1"))
	assert.True(t, base.Default)
	assert.Equal(t, defaultKey, base.Key)

	_, _, base = restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix2/path2"))
	assert.True(t, base.Default)
	assert.Equal(t, defaultKey, base.Key)

	// Restrict further the access to a specific prefix
	restricted, err = restricted.Restrict(fullPermission, SharePrefix{
		Bucket: "bucket",
		Prefix: "prefix1",
	})
	require.NoError(t, err)

	// Check that only the actions under this prefix are allowed
	assert.NoError(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action1, nil))
	assert.Error(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action2, nil))

	// Check that encAccess has a derived key for the allowed prefix instead of the default key
	_, _, base = restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix1/path1"))
	assert.False(t, base.Default)
	assert.NotEmpty(t, base.Key)
	assert.NotEqual(t, defaultKey, base.Key)
	assert.Equal(t, "prefix1", base.Encrypted.Raw())
	assert.Equal(t, "prefix1", base.Unencrypted.Raw())

	_, _, base = restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix2/path2"))
	assert.Nil(t, base)

	// Restrict further the access again only with NotAfter time.
	permission.NotAfter = time.Now().Add(1 * time.Hour)

	// Check that the access still allows only actions under the allowed prefix
	restricted, err = restricted.Restrict(permission)
	require.NoError(t, err)
	assert.NoError(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action1, nil))
	assert.Error(t, restricted.APIKey.Check(ctx, secret, macaroon.APIKeyVersionObjectLock, action2, nil))

	// Check that encAccess has not changed too
	_, _, base = restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix1/path1"))
	assert.False(t, base.Default)
	assert.NotEmpty(t, base.Key)
	assert.NotEqual(t, defaultKey, base.Key)
	assert.Equal(t, "prefix1", base.Encrypted.Raw())
	assert.Equal(t, "prefix1", base.Unencrypted.Raw())

	_, _, base = restricted.EncAccess.Store.LookupEncrypted("bucket", paths.NewEncrypted("prefix2/path2"))
	assert.Nil(t, base)
}
