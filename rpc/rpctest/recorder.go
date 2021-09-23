// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package rpctest

import (
	"sync"

	"storj.io/drpc"
)

// CallRecorder wraps drpc.Conn and record the rpc names for each calls.
// It uses an internal Mutex, therefore it's not recommended for production or
// performance critical operations.
type CallRecorder struct {
	mu    sync.Mutex
	calls []string
}

// NewCallRecorder returns with a properly initialized RPCounter.
func NewCallRecorder() CallRecorder {
	return CallRecorder{
		calls: make([]string, 0),
	}
}

// Reset deletes all the existing counters and set everything to 0.
func (r *CallRecorder) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = make([]string, 0)
}

// CountOf returns the number of calls to one specific rpc method.
func (r *CallRecorder) CountOf(rpc string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	counter := 0
	for _, c := range r.calls {
		if c == rpc {
			counter++
		}
	}
	return counter
}

// RecordCall records the fact of one rpc call.
func (r *CallRecorder) RecordCall(rpc string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, rpc)
}

// History returns the list of rpc names which called on this connection.
func (r *CallRecorder) History() []string {
	return append([]string{}, r.calls...)
}

// Attach wraps a drpc.Conn connection and returns with one where the counters are hooked in.
func (r *CallRecorder) Attach(conn drpc.Conn) drpc.Conn {
	interceptor := MessageInterceptor{
		delegate: conn,
		RequestHook: func(rpc string, message drpc.Message, err error) {
			r.RecordCall(rpc)
		},
	}
	return &interceptor
}
