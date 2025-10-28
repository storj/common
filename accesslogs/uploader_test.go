// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package accesslogs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"
	"go.uber.org/zap/zaptest"

	"storj.io/common/memory"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
)

func TestLimits(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	log := zaptest.NewLogger(t)
	defer ctx.Check(log.Sync)

	s := noopStorage{}
	u := newSequentialUploader(log, sequentialUploaderOptions{
		entryLimit:      5 * memory.KiB,
		queueLimit:      2,
		retryLimit:      1,
		shutdownTimeout: time.Second,
		uploadTimeout:   time.Minute,
	})

	for range 2 {
		require.NoError(t, u.queueUpload(s, "test", "test", testrand.Bytes(memory.KiB)))
	}
	require.ErrorIs(t, u.queueUpload(s, "test", "test", testrand.Bytes(memory.KiB)), ErrQueueLimit)
	require.ErrorIs(t, u.queueUpload(s, "test", "test", testrand.Bytes(6*memory.KiB)), ErrTooLarge)
	require.ErrorIs(t, u.queueUploadWithoutQueueLimit(s, "test", "test", testrand.Bytes(6*memory.KiB)), ErrTooLarge)
}

func TestQueueNoLimit(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	log := zaptest.NewLogger(t)
	defer ctx.Check(log.Sync)

	s := noopStorage{}
	u := newSequentialUploader(log, sequentialUploaderOptions{
		entryLimit:      5 * memory.KiB,
		queueLimit:      2,
		retryLimit:      1,
		shutdownTimeout: time.Second,
		uploadTimeout:   time.Minute,
	})
	defer ctx.Check(u.close)
	ctx.Go(u.run)

	for range 10 {
		require.NoError(t, u.queueUploadWithoutQueueLimit(s, "test", "test", testrand.Bytes(memory.KiB)))
	}
}

type errorStorage struct {
}

func (s errorStorage) Put(ctx context.Context, bucket, key string, data []byte) error {
	return errs.New("retry error")
}

func TestQueueNoLimitErroringStorage(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	log := zaptest.NewLogger(t)
	defer ctx.Check(log.Sync)

	s := errorStorage{}
	u := newSequentialUploader(log, sequentialUploaderOptions{
		entryLimit:      5 * memory.KiB,
		queueLimit:      10,
		retryLimit:      1,
		shutdownTimeout: time.Second,
		uploadTimeout:   time.Minute,
	})
	defer ctx.Check(u.close)
	ctx.Go(u.run)

	for range 10 {
		require.NoError(t, u.queueUploadWithoutQueueLimit(s, "test", "test", testrand.Bytes(memory.KiB)))
	}
}

func TestQueueErroringStorage(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	log := zaptest.NewLogger(t)
	defer ctx.Check(log.Sync)

	s := errorStorage{}
	u := newSequentialUploader(log, sequentialUploaderOptions{
		entryLimit:      5 * memory.KiB,
		queueLimit:      10,
		retryLimit:      1,
		shutdownTimeout: time.Second,
		uploadTimeout:   time.Minute,
	})
	defer ctx.Check(u.close)
	ctx.Go(u.run)

	for range 10 {
		require.NoError(t, u.queueUpload(s, "test", "test", testrand.Bytes(memory.KiB)))
	}
}

type slowStorage struct {
	delay time.Duration
}

func (s slowStorage) Put(ctx context.Context, bucket, key string, data []byte) error {
	select {
	case <-time.After(s.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestUploadTimeout(t *testing.T) {
	t.Parallel()

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	log := zaptest.NewLogger(t)
	defer ctx.Check(log.Sync)

	s := slowStorage{delay: 2 * time.Second}
	u := newSequentialUploader(log, sequentialUploaderOptions{
		entryLimit:      5 * memory.KiB,
		queueLimit:      10,
		retryLimit:      1,
		shutdownTimeout: 5 * time.Second,
		uploadTimeout:   500 * time.Millisecond,
	})
	defer ctx.Check(u.close)
	ctx.Go(u.run)

	require.NoError(t, u.queueUpload(s, "test", "test", testrand.Bytes(memory.KiB)))
}
