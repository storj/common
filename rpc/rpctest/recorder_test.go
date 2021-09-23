// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/rpc/rpctest"
	"storj.io/drpc"
)

func TestCounter(t *testing.T) {
	type helloMessage struct {
		Message string
	}

	// the original service with fake answer
	original := rpctest.NewStubConnection()
	original.RegisterHandler("func1", func(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
		out.(*helloMessage).Message = "foobar"
		return nil
	})

	callRecorder := rpctest.NewCallRecorder()
	wrapper := callRecorder.Attach(&original)

	in := helloMessage{Message: "hello"}
	out := helloMessage{}
	err := wrapper.Invoke(context.TODO(), "func1", nil, &in, &out)
	require.Nil(t, err)
	require.Equal(t, "foobar", out.Message)
	require.Equal(t, 1, callRecorder.CountOf("func1"))
}

func TestCounterAssert(t *testing.T) {
	c := rpctest.NewCallRecorder()
	c.RecordCall("first")
	c.RecordCall("first")
	c.RecordCall("second")

	// be sure we have copy, not reference
	c.History()[0] = "x"

	require.Equal(t, c.History(), []string{"first", "first", "second"})
}
