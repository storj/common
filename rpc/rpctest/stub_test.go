// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/drpc"
)

func TestStubConnection(t *testing.T) {

	// GIVEN
	type helloMessage struct {
		Message string
	}

	stub := NewStubConnection()
	defer func(stub *StubConnection) {
		require.NoError(t, stub.Close())
		// second close shouldn't cause a problem
		require.NoError(t, stub.Close())
	}(&stub)

	stub.RegisterHandler("/test/hello1", func(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
		out.(*helloMessage).Message = "replaced"
		return nil
	})

	// WHEN
	out := helloMessage{}
	err := stub.Invoke(context.Background(), "/test/hello1", nil, &helloMessage{Message: "hello"}, &out)

	// THEN
	require.Nil(t, err)
	require.Equal(t, "replaced", out.Message)

}

func TestStubConnectionAsync(t *testing.T) {
	type helloMessage struct {
		Message string
	}

	stub := NewStubConnection()
	defer func(stub *StubConnection) {
		require.NoError(t, stub.Close())
	}(&stub)

	stub.RegisterHandler("/test/hello1", func(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
		out.(*helloMessage).Message = "pong+" + in.(*helloMessage).Message
		return nil
	})

	out := helloMessage{}
	stream, err := stub.NewStream(context.Background(), "/test/hello1", nil)
	require.Nil(t, err)

	for i := range 3 {
		err = stream.MsgSend(&helloMessage{
			"ping" + strconv.Itoa(i),
		}, nil)
		require.Nil(t, err)

		err = stream.MsgRecv(&out, nil)
		require.Nil(t, err)

		require.Equal(t, "pong+ping"+strconv.Itoa(i), out.Message)
	}
}
