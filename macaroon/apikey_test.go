// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package macaroon

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"

	"storj.io/common/testcontext"
)

func TestSerializeParseRestrictAndCheck(t *testing.T) {
	ctx := context.Background()

	secret, err := NewSecret()
	require.NoError(t, err)
	key, err := NewAPIKey(secret)
	require.NoError(t, err)

	serialized := key.Serialize()
	parsedKey, err := ParseAPIKey(serialized)
	require.NoError(t, err)
	require.True(t, bytes.Equal(key.Head(), parsedKey.Head()))
	require.True(t, bytes.Equal(key.Tail(), parsedKey.Tail()))

	restricted, err := key.Restrict(WithNonce(Caveat{
		AllowedPaths: []*Caveat_Path{{
			Bucket:              []byte("a-test-bucket"),
			EncryptedPathPrefix: []byte("a-test-path"),
		}},
	}))
	require.NoError(t, err)

	serialized = restricted.Serialize()
	parsedKey, err = ParseAPIKey(serialized)
	require.NoError(t, err)
	require.True(t, bytes.Equal(key.Head(), parsedKey.Head()))
	require.False(t, bytes.Equal(key.Tail(), parsedKey.Tail()))

	now := time.Now()
	action1 := Action{
		Op:            ActionRead,
		Time:          now,
		Bucket:        []byte("a-test-bucket"),
		EncryptedPath: []byte("a-test-path"),
	}
	action2 := Action{
		Op:            ActionRead,
		Time:          now,
		Bucket:        []byte("another-test-bucket"),
		EncryptedPath: []byte("another-test-path"),
	}

	require.NoError(t, key.Check(ctx, secret, action1, nil))
	require.NoError(t, key.Check(ctx, secret, action2, nil))
	require.NoError(t, parsedKey.Check(ctx, secret, action1, nil))
	err = parsedKey.Check(ctx, secret, action2, nil)
	require.True(t, ErrUnauthorized.Has(err), err)
}

func TestRevocation(t *testing.T) {
	ctx := testcontext.New(t)

	secret, err := NewSecret()
	require.NoError(t, err)
	key, err := NewAPIKey(secret)
	require.NoError(t, err)

	now := time.Now()
	action := Action{
		Op:   ActionWrite,
		Time: now,
	}

	// Nil revoker should not revoke anything
	require.NoError(t, key.Check(ctx, secret, action, nil))

	// No error returned when nothing is revoked
	nothingRevoked := testRevoker{}
	require.NoError(t, key.Check(ctx, secret, action, nothingRevoked))

	// Check that revoked results in a ErrRevoked error
	revoked := testRevoker{revoked: true}
	require.True(t, ErrRevoked.Has(key.Check(ctx, secret, action, revoked)))

	// Check that an error while checking revocations results in an ErrRevoked err
	revokeErr := testRevoker{revoked: false, err: errs.New("some revoke err")}
	require.True(t, ErrRevoked.Has(key.Check(ctx, secret, action, revokeErr)))
}

func TestExpiration(t *testing.T) {
	ctx := context.Background()

	secret, err := NewSecret()
	require.NoError(t, err)
	key, err := NewAPIKey(secret)
	require.NoError(t, err)

	now := time.Now()
	minuteAgo := now.Add(-time.Minute)
	minuteFromNow := now.Add(time.Minute)
	twoMinutesAgo := now.Add(-2 * time.Minute)
	twoMinutesFromNow := now.Add(2 * time.Minute)

	notBeforeMinuteFromNow, err := key.Restrict(WithNonce(Caveat{
		NotBefore: &minuteFromNow,
	}))
	require.NoError(t, err)
	notAfterMinuteAgo, err := key.Restrict(WithNonce(Caveat{
		NotAfter: &minuteAgo,
	}))
	require.NoError(t, err)

	for i, test := range []struct {
		keyToTest       *APIKey
		timestampToTest time.Time
		errClass        *errs.Class
	}{
		{key, time.Time{}, &Error},
		{notBeforeMinuteFromNow, time.Time{}, &Error},
		{notAfterMinuteAgo, time.Time{}, &Error},

		{key, now, nil},
		{notBeforeMinuteFromNow, now, &ErrUnauthorized},
		{notAfterMinuteAgo, now, &ErrUnauthorized},

		{key, twoMinutesAgo, nil},
		{notBeforeMinuteFromNow, twoMinutesAgo, &ErrUnauthorized},
		{notAfterMinuteAgo, twoMinutesAgo, nil},

		{key, twoMinutesFromNow, nil},
		{notBeforeMinuteFromNow, twoMinutesFromNow, nil},
		{notAfterMinuteAgo, twoMinutesFromNow, &ErrUnauthorized},
	} {
		err := test.keyToTest.Check(ctx, secret, Action{
			Op:   ActionRead,
			Time: test.timestampToTest,
		}, nil)
		if test.errClass == nil {
			require.NoError(t, err, fmt.Sprintf("test #%d", i+1))
		} else {
			require.False(t, !test.errClass.Has(err), fmt.Sprintf("test #%d", i+1))
		}
	}
}

func TestGetAllowedBuckets(t *testing.T) {
	ctx := context.Background()

	secret, err := NewSecret()
	require.NoError(t, err)
	key, err := NewAPIKey(secret)
	require.NoError(t, err)

	restricted, err := key.Restrict(WithNonce(Caveat{
		AllowedPaths: []*Caveat_Path{{Bucket: []byte("test1")}, {Bucket: []byte("test2")}},
	}))
	require.NoError(t, err)

	restricted, err = restricted.Restrict(WithNonce(Caveat{
		AllowedPaths: []*Caveat_Path{{Bucket: []byte("test1")}, {Bucket: []byte("test3")}},
	}))
	require.NoError(t, err)

	now := time.Now()
	action := Action{
		Op:   ActionRead,
		Time: now,
	}

	allowed, err := restricted.GetAllowedBuckets(ctx, action)
	require.NoError(t, err)
	require.Equal(t, allowed, AllowedBuckets{
		All:     false,
		Buckets: map[string]struct{}{"test1": {}},
	})

	restricted, err = restricted.Restrict(WithNonce(Caveat{DisallowWrites: true}))
	require.NoError(t, err)

	allowed, err = restricted.GetAllowedBuckets(ctx, action)
	require.NoError(t, err)
	require.Equal(t, allowed, AllowedBuckets{
		All:     false,
		Buckets: map[string]struct{}{"test1": {}},
	})
}

func TestNonce(t *testing.T) {
	secret, err := NewSecret()
	require.NoError(t, err)

	key1, err := NewAPIKey(secret)
	require.NoError(t, err)

	key2, err := ParseAPIKey(key1.Serialize())
	require.NoError(t, err)

	// Key 1 and 2 should be exactly equal.
	require.True(t, bytes.Equal(key1.Head(), key2.Head()))
	require.True(t, bytes.Equal(key1.Tail(), key2.Tail()))

	caveat := Caveat{
		DisallowReads: true,
	}

	t.Run("WithoutNonce", func(t *testing.T) {
		key1r, err := key1.Restrict(caveat)
		require.NoError(t, err)

		key2r, err := key2.Restrict(caveat)
		require.NoError(t, err)

		// Key 1 and 2 should be exactly equal when the caveat does not have a
		// nonce.
		require.True(t, bytes.Equal(key1r.Head(), key2r.Head()))
		require.True(t, bytes.Equal(key1r.Tail(), key2r.Tail()))
	})

	t.Run("WithNonce", func(t *testing.T) {
		key1r, err := key1.Restrict(WithNonce(caveat))
		require.NoError(t, err)

		key2r, err := key2.Restrict(WithNonce(caveat))
		require.NoError(t, err)

		// Key 1 and 2 should share the same head, but have different
		// tails when the caveats have a nonce.
		require.True(t, bytes.Equal(key1r.Head(), key2r.Head()))
		require.False(t, bytes.Equal(key1r.Tail(), key2r.Tail()))
	})
}

func BenchmarkAPIKey_Check(b *testing.B) {
	ctx := context.Background()

	secret, err := NewSecret()
	require.NoError(b, err)

	key, err := NewAPIKey(secret)
	require.NoError(b, err)

	now := time.Now()
	denyBefore := now.Add(-10 * time.Hour)
	denyAfter := now.Add(10 * time.Hour)

	key2, err := key.Restrict(WithNonce(Caveat{
		NotBefore: &denyBefore,
		AllowedPaths: []*Caveat_Path{
			{Bucket: []byte("test1")},
			{Bucket: []byte("test3")},
		},
	}))
	require.NoError(b, err)

	key3, err := key2.Restrict(WithNonce(Caveat{
		NotAfter: &denyAfter,
	}))
	require.NoError(b, err)

	b.ResetTimer()

	revoker := &testRevoker{}
	for i := 0; i < b.N; i++ {
		err := key3.Check(ctx, secret, Action{Bucket: []byte("test1"), Op: ActionRead, Time: now}, revoker)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type testRevoker struct {
	revoked bool
	err     error
}

func (tr testRevoker) Check(ctx context.Context, tails [][]byte) (bool, error) {
	return tr.revoked, tr.err
}
