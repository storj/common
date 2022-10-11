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
	require.False(t, Has(ctx, "first"))

	ctx = With(ctx, "first")
	require.True(t, Has(ctx, "first"))
	require.False(t, Has(ctx, "second"))

	ctx = With(ctx, "second")
	require.True(t, Has(ctx, "first"))
	require.True(t, Has(ctx, "second"))

	exps := Get(ctx)
	sort.Strings(exps)
	require.Equal(t, []string{"first", "second"}, exps)

}
