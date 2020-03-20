// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// Package rpctracing implements tracing for rpc.
package rpctracing

import (
	"context"
	"strconv"

	"github.com/spacemonkeygo/monkit/v3"

	"storj.io/drpc"
	"storj.io/drpc/drpcmetadata"
	"storj.io/drpc/drpcmux"
)

type streamWrapper struct {
	drpc.Stream
	ctx context.Context
}

func (s *streamWrapper) Context() context.Context { return s.ctx }

type handlerFunc func(traceId *int64, parentId *int64) (trace *monkit.Trace, spanId int64)

var defaultHandlerFunc = func(traceId *int64, parentId *int64) (*monkit.Trace, int64) {
	return monkit.NewTrace(*traceId), monkit.NewId()
}

// Handler implements drpc handler interface and takes in a callback function.
type Handler struct {
	mux *drpcmux.Mux
	cb  handlerFunc
}

// NewHandler returns a new instance of Handler.
func NewHandler(mux *drpcmux.Mux, cb handlerFunc) *Handler {
	if cb == nil {
		cb = defaultHandlerFunc
	}
	return &Handler{
		mux: mux,
		cb:  cb,
	}
}

// HandleRPC adds tracing metadata onto server stream.
func (handler *Handler) HandleRPC(stream drpc.Stream, rpc string) (err error) {
	streamCtx := stream.Context()
	metadata, ok := drpcmetadata.Get(streamCtx)
	if ok {
		parentID, err := strconv.ParseInt(metadata[ParentID], 10, 64)
		if err != nil {
			return handler.mux.HandleRPC(stream, rpc)
		}

		traceID, err := strconv.ParseInt(metadata[TraceID], 10, 64)
		if err != nil {
			return handler.mux.HandleRPC(stream, rpc)
		}
		trace, spanID := handler.cb(&traceID, &parentID)
		defer mon.FuncNamed(rpc).RemoteTrace(&streamCtx, spanID, trace)(&err)
	}

	return handler.mux.HandleRPC(&streamWrapper{Stream: stream, ctx: streamCtx}, rpc)
}
