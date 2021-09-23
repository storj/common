// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/drpc"
)

type message struct {
	content string
}

func TestMessageInterceptor(t *testing.T) {
	requestCounter := 0
	responseCounter := 0

	// the original service with fake answer
	original := NewStubConnection()
	original.RegisterHandler("test", func(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
		out.(*message).content = "foobar"
		return nil
	})

	// the wrapper, counts the invocations
	wrapper := NewMessageInterceptor(&original)
	wrapper.RequestHook = func(rpc string, message drpc.Message, err error) {
		requestCounter++
	}
	wrapper.ResponseHook = func(rpc string, message drpc.Message, err error) {
		responseCounter++
	}

	in := message{content: "hello"}
	out := message{}

	// sync call
	err := wrapper.Invoke(context.TODO(), "test", nil, &in, &out)
	require.Nil(t, err)
	require.Equal(t, "foobar", out.content)

	out = message{}

	// async call
	stream, err := wrapper.NewStream(context.TODO(), "test", nil)
	require.Nil(t, err)

	// send it
	err = stream.MsgSend(&in, nil)
	require.Nil(t, err)

	// wait for the response
	err = stream.MsgRecv(&out, nil)
	require.Nil(t, err)
	require.Equal(t, "foobar", out.content)

	require.Equal(t, 2, requestCounter)
	require.Equal(t, 2, responseCounter)

}
