// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package sync2

import (
	"context"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
)

func requireBlocked(t *testing.T, ch chan struct{}) {
	select {
	case <-ch:
		t.Fatal("channel read or close when it should not have been")
	default:
	}
}

func TestReceiverClosableChan_Basic(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	ch := MakeReceiverClosableChan[int](3)

	require.True(t, ch.BlockingSend(1))
	require.True(t, ch.BlockingSend(2))

	v, err := ch.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, v, 1)
	v, err = ch.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, v, 2)

	require.True(t, ch.BlockingSend(3))
	require.True(t, ch.BlockingSend(4))
	require.True(t, ch.BlockingSend(5))

	v, err = ch.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, v, 3)

	vs := ch.StopReceiving()
	require.Equal(t, vs, []int{4, 5})

	require.False(t, ch.BlockingSend(6))
}

func TestReceiverClosableChan_BlockingSend(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	ch := MakeReceiverClosableChan[int](3)

	require.True(t, ch.BlockingSend(1))
	require.True(t, ch.BlockingSend(2))
	require.True(t, ch.BlockingSend(3))

	sending := make(chan struct{})
	sent := make(chan struct{})
	sentWithoutRace := false

	ctx.Go(func() error {
		close(sending)
		require.True(t, ch.BlockingSend(4))
		close(sent)
		sentWithoutRace = true
		return nil
	})

	<-sending
	for range 10 {
		// make sure the send is blocked
		runtime.Gosched()
	}
	requireBlocked(t, sent)
	require.False(t, sentWithoutRace)

	v, err := ch.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, v, 1)
	<-sent
}

func TestReceiverClosableChan_UnableToSend(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	ch := MakeReceiverClosableChan[int](3)

	require.True(t, ch.BlockingSend(1))
	require.True(t, ch.BlockingSend(2))
	require.True(t, ch.BlockingSend(3))

	sending := make(chan struct{})
	sent := make(chan struct{})
	sentWithoutRace := false

	ctx.Go(func() error {
		close(sending)
		require.False(t, ch.BlockingSend(4))
		close(sent)
		sentWithoutRace = true
		return nil
	})

	<-sending
	for range 10 {
		// make sure the send is blocked
		runtime.Gosched()
	}
	requireBlocked(t, sent)
	require.False(t, sentWithoutRace)

	vs := ch.StopReceiving()
	require.Equal(t, vs, []int{1, 2, 3})
	<-sent
}

func TestReceiverClosableChan_BlockingReceive(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	ch := MakeReceiverClosableChan[int](3)

	receiving := make(chan struct{})
	received := make(chan struct{})
	receivedWithoutRace := false

	ctx.Go(func() error {
		close(receiving)
		ctx := context.Background()
		v, err := ch.Receive(ctx)
		require.NoError(t, err)
		require.Equal(t, v, 1)
		close(received)
		receivedWithoutRace = true
		return nil
	})

	<-receiving
	for range 10 {
		// make sure the receive is blocked
		runtime.Gosched()
	}
	requireBlocked(t, received)
	require.False(t, receivedWithoutRace)

	require.True(t, ch.BlockingSend(1))
	<-received
}

func TestReceiverClosableChan_ContextCanceled(t *testing.T) {
	t.Parallel()
	ctx := testcontext.New(t)

	ch := MakeReceiverClosableChan[int](3)

	receiving := make(chan struct{})
	received := make(chan struct{})
	receivedWithoutRace := false
	cancelCtx, cancel := context.WithCancel(ctx)

	ctx.Go(func() error {
		close(receiving)
		_, err := ch.Receive(cancelCtx)
		require.ErrorIs(t, err, context.Canceled)
		close(received)
		receivedWithoutRace = true
		return nil
	})

	<-receiving
	for range 10 {
		// make sure the receive is blocked
		runtime.Gosched()
	}
	requireBlocked(t, received)
	require.False(t, receivedWithoutRace)

	cancel()

	toSend := 3
	var expected []int
	for i := range toSend {
		require.True(t, ch.BlockingSend(i))
		expected = append(expected, i)
	}

	<-received

	require.Equal(t, ch.StopReceiving(), expected)
}
