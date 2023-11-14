// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package signing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/storj"
	"storj.io/common/testcontext"
)

func TestSigningWhenDisabled(t *testing.T) {
	ctx := withSignaturesDisabledForTest(testcontext.New(t))
	privateKey := storj.PiecePrivateKey{}

	t.Run("SignOrderLimit", func(t *testing.T) {
		unsigned := &pb.OrderLimit{}
		signed, err := SignOrderLimit(ctx, panicSigner{}, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.OrderLimit{SatelliteSignature: disabledSignature}, signed)
	})
	t.Run("SignUplinkOrder", func(t *testing.T) {
		unsigned := &pb.Order{}
		signed, err := SignUplinkOrder(ctx, privateKey, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.Order{UplinkSignature: disabledSignature}, signed)
	})
	t.Run("SignPieceHash", func(t *testing.T) {
		unsigned := &pb.PieceHash{}
		signed, err := SignPieceHash(ctx, panicSigner{}, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.PieceHash{Signature: disabledSignature}, signed)
	})
	t.Run("SignUplinkPieceHash", func(t *testing.T) {
		unsigned := &pb.PieceHash{}
		signed, err := SignUplinkPieceHash(ctx, privateKey, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.PieceHash{Signature: disabledSignature}, signed)
	})
	t.Run("SignExitCompleted", func(t *testing.T) {
		unsigned := &pb.ExitCompleted{}
		signed, err := SignExitCompleted(ctx, panicSigner{}, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.ExitCompleted{ExitCompleteSignature: disabledSignature}, signed)
	})
	t.Run("SignExitFailed", func(t *testing.T) {
		unsigned := &pb.ExitFailed{}
		signed, err := SignExitFailed(ctx, panicSigner{}, unsigned)
		require.NoError(t, err)
		assert.NotSame(t, unsigned, signed)
		assert.Equal(t, &pb.ExitFailed{ExitFailureSignature: disabledSignature}, signed)
	})
}

func TestVerificationWhenDisabled(t *testing.T) {
	ctx := withSignaturesDisabledForTest(testcontext.New(t))
	publicKey := storj.PiecePublicKey{}
	badSignature := []byte("BADBADBAD")

	t.Run("VerifyOrderLimitSignature", func(t *testing.T) {
		err := VerifyOrderLimitSignature(ctx, panicSignee{}, &pb.OrderLimit{SatelliteSignature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyOrderSignature", func(t *testing.T) {
		err := VerifyOrderSignature(ctx, panicSignee{}, &pb.Order{UplinkSignature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyUplinkOrderSignature", func(t *testing.T) {
		err := VerifyUplinkOrderSignature(ctx, publicKey, &pb.Order{UplinkSignature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyPieceHashSignature", func(t *testing.T) {
		err := VerifyPieceHashSignature(ctx, panicSignee{}, &pb.PieceHash{Signature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyUplinkPieceHashSignature", func(t *testing.T) {
		err := VerifyUplinkPieceHashSignature(ctx, publicKey, &pb.PieceHash{Signature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyExitCompleted", func(t *testing.T) {
		err := VerifyExitCompleted(ctx, panicSignee{}, &pb.ExitCompleted{ExitCompleteSignature: badSignature})
		require.NoError(t, err)
	})

	t.Run("VerifyExitFailed", func(t *testing.T) {
		err := VerifyExitFailed(ctx, panicSignee{}, &pb.ExitFailed{ExitFailureSignature: badSignature})
		require.NoError(t, err)
	})

}

type panicSigner struct{}

func (panicSigner) ID() storj.NodeID {
	panic("should not be called")
}

func (panicSigner) HashAndSign(ctx context.Context, data []byte) ([]byte, error) {
	panic("should not be called")
}

func (panicSigner) HashAndVerifySignature(ctx context.Context, data []byte, signature []byte) error {
	panic("should not be called")
}

func (panicSigner) SignHMACSHA256(ctx context.Context, data []byte) ([]byte, error) {
	panic("should not be called")
}

func (panicSigner) VerifyHMACSHA256(ctx context.Context, data []byte, signature []byte) error {
	panic("should not be called")
}

type panicSignee struct{}

func (panicSignee) ID() storj.NodeID {
	panic("should not be called")
}

func (panicSignee) HashAndVerifySignature(ctx context.Context, data []byte, signature []byte) error {
	panic("should not be called")
}
