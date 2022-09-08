// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// Package experiment implements feature flag propagation.
package experiment

import (
	"context"
)

type key int

const (
	contextKey key    = iota
	drpcKey    string = "experiment"
)

// WithExperiment registers the feature flag of an ongoing experiment.
func WithExperiment(ctx context.Context, experiment string) context.Context {
	if experiment != "" {
		return context.WithValue(ctx, contextKey, experiment)
	}
	return ctx
}

// GetExperiment returns the registered feature flag.
func GetExperiment(ctx context.Context) string {
	value := ctx.Value(contextKey)
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}
