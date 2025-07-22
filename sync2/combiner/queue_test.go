// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information

package combiner_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2/combiner"
)

func TestLimitedQueue(t *testing.T) {
	t.Parallel()

	const jobsPerQueue = 3
	q := combiner.NewQueue[int](jobsPerQueue)

	var g errgroup.Group
	for i := range jobsPerQueue {
		g.Go(func() error {
			ok := q.TryPush(i)
			if !ok {
				return errors.New("expected to push")
			}
			return nil
		})
	}
	require.NoError(t, g.Wait())

	require.False(t, q.TryPush(-1))

	values, ok := q.PopAll()
	require.True(t, ok)

	sort.Ints(values)
	require.EqualValues(t, []int{0, 1, 2}, values)
}

func TestLimitedQueue_Unbounded(t *testing.T) {
	t.Parallel()

	const tryPush = 100
	q := combiner.NewQueue[int](-1)

	var g errgroup.Group
	expect := []int{}
	for i := range tryPush {
		expect = append(expect, i)
		g.Go(func() error {
			ok := q.TryPush(i)
			if !ok {
				return errors.New("expected to push")
			}
			return nil
		})
	}
	require.NoError(t, g.Wait())

	values, ok := q.PopAll()
	require.True(t, ok)

	sort.Ints(values)
	require.EqualValues(t, expect, values)
}
