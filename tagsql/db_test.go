// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package tagsql_test

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"storj.io/common/dbutil/cockroachutil"
	"storj.io/common/dbutil/pgtest"
	"storj.io/common/dbutil/pgutil"
	"storj.io/common/tagsql"
	"storj.io/common/testcontext"
)

func run(t *testing.T, fn func(*testcontext.Context, *testing.T, tagsql.DB, tagsql.ContextSupport)) {
	t.Helper()

	t.Run("mattn-sqlite3", func(t *testing.T) {
		ctx := testcontext.New(t)
		defer ctx.Cleanup()

		db, err := tagsql.Open(ctx, "sqlite3", ":memory:")
		if err != nil {
			t.Fatal(err)
		}
		defer ctx.Check(db.Close)

		fn(ctx, t, db, tagsql.SupportBasic)
	})

	t.Run("jackc-pgx-postgres", func(t *testing.T) {
		connstr := pgtest.PickPostgres(t)

		ctx := testcontext.New(t)
		defer ctx.Cleanup()

		db, err := pgutil.OpenUnique(ctx, connstr, "detect")
		require.NoError(t, err)
		defer ctx.Check(db.Close)

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(100)

		fn(ctx, t, db.DB, tagsql.SupportAll)
	})

	t.Run("jackc-pgx-cockroach", func(t *testing.T) {
		connstr := pgtest.PickCockroach(t)

		ctx := testcontext.New(t)
		defer ctx.Cleanup()

		db, err := cockroachutil.OpenUnique(ctx, connstr, "detect")
		require.NoError(t, err)
		defer ctx.Check(db.Close)

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(100)

		fn(ctx, t, db.DB, tagsql.SupportAll)
	})
}
