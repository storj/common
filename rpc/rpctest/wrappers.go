// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"context"

	"storj.io/drpc"
)

// MessageHook is a function which may be called before and after an rpc call.
type MessageHook func(rpc string, message drpc.Message, err error)

// MessageInterceptor is a drpc.Conn which wraps an original connection with more functionality.
type MessageInterceptor struct {
	delegate     drpc.Conn
	RequestHook  MessageHook
	ResponseHook MessageHook
}

// NewMessageInterceptor creates a MessageInterceptor, a connection which delegates all the call to the specific drpc.Conn.
func NewMessageInterceptor(conn drpc.Conn) MessageInterceptor {
	return MessageInterceptor{
		delegate: conn,
	}
}

// Close closes underlying dprc connection.
func (l *MessageInterceptor) Close() error {
	return l.delegate.Close()
}

// Closed returns a channel that is closed if the underlying connection is definitely closed.
func (l *MessageInterceptor) Closed() <-chan struct{} {
	return l.delegate.Closed()
}

// Invoke the underlying connection but call the RequestHook/ResponseHook before and after.
// When the Invoker is set it will be invoked instead of the original connection.
func (l *MessageInterceptor) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in, out drpc.Message) error {
	var err error
	if l.RequestHook != nil {
		l.RequestHook(rpc, in, nil)
	}

	err = l.delegate.Invoke(ctx, rpc, enc, in, out)

	if l.ResponseHook != nil {
		l.ResponseHook(rpc, out, err)
	}
	return err
}

// NewStream creates a new wrapped stream.
func (l *MessageInterceptor) NewStream(ctx context.Context, rpc string, enc drpc.Encoding) (drpc.Stream, error) {
	stream, err := l.delegate.NewStream(ctx, rpc, enc)
	if err != nil {
		return stream, err
	}
	return &interceptedStream{
		delegate:     stream,
		rpc:          rpc,
		requestHook:  l.RequestHook,
		responseHook: l.ResponseHook,
	}, nil
}

type interceptedStream struct {
	delegate     drpc.Stream
	requestHook  MessageHook
	responseHook MessageHook
	rpc          string
}

// Context returns the context from the underlying stream.
func (d *interceptedStream) Context() context.Context {
	return d.delegate.Context()
}

// MsgSend sends the Message to the underlying remote OR calls the configured Invoker if defined.
// In both cases the RequestHook is called before.
func (d *interceptedStream) MsgSend(msg drpc.Message, enc drpc.Encoding) error {
	var err error
	if d.requestHook != nil {
		d.requestHook(d.rpc, msg, err)
	}

	err = d.delegate.MsgSend(msg, enc)
	return err
}

// MsgRecv receives a Message from the underlying wrapped remote.
// The configured responseHook is executed before return.
func (d *interceptedStream) MsgRecv(msg drpc.Message, enc drpc.Encoding) error {
	err := d.delegate.MsgRecv(msg, enc)
	if d.responseHook != nil {
		d.responseHook(d.rpc, msg, err)
	}
	return err
}

// CloseSend signals to the remote that we will no longer send any messages.
func (d *interceptedStream) CloseSend() error {
	return d.delegate.CloseSend()
}

// Close closes the stream.
func (d *interceptedStream) Close() error {
	return d.delegate.Close()
}
