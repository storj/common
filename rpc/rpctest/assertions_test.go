// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/common/rpc/rpcstatus"
)

func TestAssertCode(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("whatever"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		assert.True(t, AssertCode(fake, err, rpcstatus.Internal))
		assert.False(t, fake.failed)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		assert.False(t, AssertCode(fake, err, rpcstatus.InvalidArgument))
		assert.True(t, fake.failed)
	})
}

func TestRequireCode(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("whatever"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		RequireCode(fake, err, rpcstatus.Internal)
		assert.False(t, fake.failed)
		assert.False(t, fake.failNow)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		RequireCode(fake, err, rpcstatus.InvalidArgument)
		assert.True(t, fake.failed)
		assert.True(t, fake.failNow)
	})
}

func TestAssertStatus(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("you shall pass"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		assert.True(t, AssertStatus(fake, err, rpcstatus.Internal, "you shall pass"))
		assert.False(t, fake.failed)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		assert.False(t, AssertStatus(fake, err, rpcstatus.InvalidArgument, "you shall pass"))
		assert.True(t, fake.failed)
	})

	t.Run("fail with unexpected cause", func(t *testing.T) {
		fake := new(fakeTB)
		assert.False(t, AssertStatus(fake, err, rpcstatus.Internal, "you shall not pass"))
		assert.True(t, fake.failed)
	})
}

func TestRequireStatus(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("you shall pass"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatus(fake, err, rpcstatus.Internal, "you shall pass")
		assert.False(t, fake.failed)
		assert.False(t, fake.failNow)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatus(fake, err, rpcstatus.InvalidArgument, "you shall pass")
		assert.True(t, fake.failed)
		assert.True(t, fake.failNow)
	})

	t.Run("fail with unexpected cause", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatus(fake, err, rpcstatus.Internal, "you shall not pass")
		assert.True(t, fake.failed)
		assert.True(t, fake.failNow)
	})
}

func TestAssertStatusContains(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("you shall pass"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		assert.True(t, AssertStatusContains(fake, err, rpcstatus.Internal, "shall pass"))
		assert.False(t, fake.failed)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		assert.False(t, AssertStatusContains(fake, err, rpcstatus.InvalidArgument, "shall pass"))
		assert.True(t, fake.failed)
	})

	t.Run("fail with unexpected cause", func(t *testing.T) {
		fake := new(fakeTB)
		assert.False(t, AssertStatusContains(fake, err, rpcstatus.Internal, "shall not pass"))
		assert.True(t, fake.failed)
	})
}

func TestRequireStatusContains(t *testing.T) {
	err := rpcstatus.Wrap(rpcstatus.Internal, errors.New("you shall pass"))

	t.Run("pass", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatusContains(fake, err, rpcstatus.Internal, "shall pass")
		assert.False(t, fake.failed)
		assert.False(t, fake.failNow)
	})

	t.Run("fail with unexpected code", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatusContains(fake, err, rpcstatus.InvalidArgument, "shall pass")
		assert.True(t, fake.failed)
		assert.True(t, fake.failNow)
	})

	t.Run("fail with unexpected cause", func(t *testing.T) {
		fake := new(fakeTB)
		RequireStatusContains(fake, err, rpcstatus.Internal, "shall not pass")
		assert.True(t, fake.failed)
		assert.True(t, fake.failNow)
	})
}

type fakeTB struct {
	// TB embedded so we don't have to provide all the methods. If we miss
	// one, this will cause a panic, which is ok.
	testing.TB

	failed  bool
	helper  bool
	failNow bool
}

func (tb *fakeTB) Name() string {
	return "fake"
}

func (tb *fakeTB) Errorf(format string, args ...any) {
	tb.failed = true
}

func (tb *fakeTB) Helper() {
	tb.helper = true
}

func (tb *fakeTB) FailNow() {
	tb.failNow = true
}
