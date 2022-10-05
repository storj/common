// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package experiment

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetExperimentInContext(t *testing.T) {
	ctx := context.Background()
	require.False(t, HasExperiment(ctx, "first"))

	ctx = WithExperiment(ctx, "first")
	require.True(t, HasExperiment(ctx, "first"))
	require.False(t, HasExperiment(ctx, "second"))

	ctx = WithExperiment(ctx, "second")
	require.True(t, HasExperiment(ctx, "first"))
	require.True(t, HasExperiment(ctx, "second"))

	exps := GetExperiment(ctx)
	sort.Strings(exps)
	require.Equal(t, []string{"first", "second"}, exps)

}
