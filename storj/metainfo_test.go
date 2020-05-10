// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package storj_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
)

func TestListOptions(t *testing.T) {
	opts := storj.ListOptions{
		Prefix:    "alpha/",
		Cursor:    "a",
		Delimiter: '/',
		Recursive: true,
		Direction: storj.After,
		Limit:     30,
	}

	list := storj.ObjectList{
		Bucket: "hello",
		Prefix: "alpha/",
		More:   true,
		Items: []storj.Object{
			{Path: "alpha/xyz"},
		},
	}

	newopts := opts.NextPage(list)
	require.Equal(t, storj.ListOptions{
		Prefix:    "alpha/",
		Cursor:    "alpha/xyz",
		Delimiter: '/',
		Recursive: true,
		Direction: storj.After,
		Limit:     30,
	}, newopts)
}
