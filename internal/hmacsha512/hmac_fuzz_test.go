// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.18
// +build go1.18

package hmacsha512_test

import (
	"crypto/hmac"
	"crypto/sha512"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/internal/hmacsha512"
)

func FuzzPartial(f *testing.F) {
	pid := PieceID{}
	f.Add([]byte{}, []byte{})
	f.Add(pid[:], []byte{1, 2, 3, 4})

	f.Fuzz(func(t *testing.T, key []byte, data []byte) {
		var local hmacsha512.Partial
		local.Init(key)
		local.Write(data)
		actual1 := local.SumAndReset()
		local.Write(data)
		actual2 := local.SumAndReset()

		slow := hmac.New(sha512.New, key)
		_, _ = slow.Write(data)
		expected := slow.Sum(nil)

		require.Equal(t, expected, actual1[:])
		require.Equal(t, expected, actual2[:])
	})
}
