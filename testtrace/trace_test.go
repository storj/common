// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package testtrace_test

import (
	"context"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/testtrace"
)

func TestLabels(t *testing.T) {
	const label = "LABEL IS VISIBLE"

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	pprof.Do(ctx, pprof.Labels("name", label), func(c context.Context) {
		for range 3 {
			go func() { <-c.Done() }()
		}
	})

	time.Sleep(time.Millisecond)

	trace, err := testtrace.Summary()
	require.NoError(t, err)
	t.Log("\n" + trace)

	require.Contains(t, trace, label)
}

func TestFilter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	pprof.Do(ctx, pprof.Labels("name", "alpha", "user", "ALPHA"), func(c context.Context) {
		go func() { <-c.Done() }()
	})
	pprof.Do(ctx, pprof.Labels("name", "beta", "user", "BETA"), func(c context.Context) {
		go func() { <-c.Done() }()
	})

	time.Sleep(time.Millisecond)

	trace, err := testtrace.Summary("name", "alpha")
	require.NoError(t, err)
	t.Log("\n" + trace)

	require.Contains(t, trace, "ALPHA")
	require.NotContains(t, trace, "alpha")
	require.NotContains(t, trace, "BETA")
}
