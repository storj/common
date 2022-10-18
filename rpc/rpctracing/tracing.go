// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// Package rpctracing implements tracing for rpc.
package rpctracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"strconv"

	"github.com/spacemonkeygo/monkit/v3"

	"storj.io/common/rpc/rpcpool"
	"storj.io/drpc"
	"storj.io/drpc/drpcmetadata"
)

// TracingWrapper wraps a drpc.Conn with tracing information.
type TracingWrapper struct {
	rpcpool.Conn
}

// NewTracingWrapper creates a new instance of the wrapper.
func NewTracingWrapper(conn rpcpool.Conn) *TracingWrapper {
	return &TracingWrapper{
		Conn: conn,
	}
}

// Invoke implements drpc.Conn's Invoke method with tracing information injected into the context.
func (c *TracingWrapper) Invoke(ctx context.Context, rpc string, enc drpc.Encoding, in drpc.Message, out drpc.Message) (err error) {
	return c.Conn.Invoke(c.trace(ctx, rpc), rpc, enc, in, out)
}

// NewStream implements drpc.Conn's NewStream method with tracing information injected into the context.
func (c *TracingWrapper) NewStream(ctx context.Context, rpc string, enc drpc.Encoding) (_ drpc.Stream, err error) {
	return c.Conn.NewStream(c.trace(ctx, rpc), rpc, enc)
}

// trace injects tracing related information into the context.
func (c *TracingWrapper) trace(ctx context.Context, rpc string) context.Context {
	ctx, otelSpan := otel.Tracer("").Start(ctx, rpc)
	defer otelSpan.End()
	metadata := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(metadata))

	monkitSpan := monkit.SpanFromCtx(ctx)
	if monkitSpan == nil {
		return ctx
	}

	sampled, exist := monkitSpan.Trace().Get(Sampled).(bool)
	if !exist || !sampled {
		return ctx
	}

	metadata[TraceID] = strconv.FormatInt(monkitSpan.Trace().Id(), 10)
	metadata[ParentID] = strconv.FormatInt(monkitSpan.Id(), 10)
	metadata[Sampled] = strconv.FormatBool(sampled)

	return drpcmetadata.AddPairs(ctx, metadata)

}
