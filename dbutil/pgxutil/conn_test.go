// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package pgxutil_test

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"storj.io/common/testcontext"
	"storj.io/common/dbutil/pgtest"
	"storj.io/common/dbutil/pgxutil"
	"storj.io/common/dbutil/tempdb"
)

func TestConn(t *testing.T) {
	pgtest.Run(t, func(ctx *testcontext.Context, t *testing.T, connstr string) {
		db, err := tempdb.OpenUnique(ctx, connstr, "pgutil-query")
		require.NoError(t, err)
		defer ctx.Check(db.Close)

		require.NoError(t,
			pgxutil.Conn(ctx, db.DB, func(conn *pgx.Conn) error {
				return nil
			}))

		require.Error(t,
			pgxutil.Conn(ctx, db.DB, func(conn *pgx.Conn) error {
				return errors.New("xyz")
			}))
	})
}
