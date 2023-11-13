// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"context"
	"fmt"
	"sync"

	"storj.io/drpc"
)

// MessageHandler is the logic called with the rpc message instead of the original backend.
type MessageHandler func(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error

// StubConnection can wrap existing drpc.Conn and replaces.
type StubConnection struct {
	invokers map[string]MessageHandler
	closed   chan struct{}
	close    sync.Once
}

// NewStubConnection create a properly initialized StubConnection.
func NewStubConnection() StubConnection {
	return StubConnection{
		invokers: make(map[string]MessageHandler),
		closed:   make(chan struct{}),
	}
}

// RegisterHandler saves the handler for specific rpc calls.
func (s *StubConnection) RegisterHandler(rpc string, invoker MessageHandler) {
	s.invokers[rpc] = invoker
}

// Close simulates the close call (noop).
func (s *StubConnection) Close() error {
	s.close.Do(func() {
		close(s.closed)
	})

	return nil
}

// Closed returns a channel that is closed if the underlying connection is definitely closed.
func (s *StubConnection) Closed() <-chan struct{} {
	return s.closed
}

// Unblocked returns a closed channel.
func (s *StubConnection) Unblocked() <-chan struct{} {
	x := make(chan struct{})
	close(x)
	return x
}

// Transport returns a nil transport.
func (s *StubConnection) Transport() drpc.Transport {
	return nil
}

// Invoke the underlying connection but call the RequestHook/ResponseHook before and after.
// When the Invoker is set it will be invoked instead of the original connection.
func (s *StubConnection) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
	invoker, found := s.invokers[rpc]
	if !found {
		return fmt.Errorf("no invoker is registered for rpc type %s", rpc)
	}
	return invoker(ctx, rpc, enc, in, out)
}

// NewStream creates a new wrapped stream.
func (s *StubConnection) NewStream(ctx context.Context, rpc string, enc drpc.Encoding) (drpc.Stream, error) {
	return &stubStream{
		invokers: s.invokers,
		rpc:      rpc,
		ctx:      ctx,
		messages: make(chan drpc.Message, 20),
	}, nil
}

type stubStream struct {
	invokers map[string]MessageHandler
	rpc      string
	ctx      context.Context
	messages chan drpc.Message
}

// Context returns the context from the underlying stream.
func (s *stubStream) Context() context.Context {
	return s.ctx
}

// MsgSend sends the Message to the underlying remote OR calls the configured Invoker if defined.
// In both cases the RequestHook is called before.
func (s *stubStream) MsgSend(msg drpc.Message, enc drpc.Encoding) error {
	s.messages <- msg
	return nil
}

// MsgRecv simulates receiving a message.
func (s *stubStream) MsgRecv(msg drpc.Message, enc drpc.Encoding) error {
	in := <-s.messages
	invoker, found := s.invokers[s.rpc]
	if !found {
		return fmt.Errorf("no invoker is registered for rpc type %s", s.rpc)
	}
	return invoker(s.ctx, s.rpc, enc, in, msg)
}

// CloseSend signals to the remote that we will no longer send any messages (no action on stub).
func (s *stubStream) CloseSend() error {
	return nil
}

// Close simulates closes the stream.
func (s *stubStream) Close() error {
	close(s.messages)
	return nil
}
