// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package errs2_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/errs2"
)

func TestGroup(t *testing.T) {
	group := errs2.Group{}
	group.Go(func() error {
		return errors.New("first")
	})
	group.Go(func() error {
		return nil
	})
	group.Go(func() error {
		return errors.New("second")
	})
	group.Go(func() error {
		return errors.New("third")
	})

	allErrors := group.Wait()
	require.Len(t, allErrors, 3)
}
