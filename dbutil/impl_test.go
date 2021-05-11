// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package dbutil_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"storj.io/private/dbutil"
)

func TestAsOfSystemTime(t *testing.T) {
	tests := []struct {
		impl dbutil.Implementation
		time time.Time
		exp  string
	}{
		{impl: dbutil.Unknown, time: time.Time{}, exp: ""},
		{impl: dbutil.Postgres, time: time.Time{}, exp: ""},
		{impl: dbutil.Cockroach, time: time.Time{}, exp: ""},
		{impl: dbutil.Bolt, time: time.Time{}, exp: ""},
		{impl: dbutil.Redis, time: time.Time{}, exp: ""},
		{impl: dbutil.SQLite3, time: time.Time{}, exp: ""},

		{impl: dbutil.Unknown, time: time.Unix(0, 1620721781789035200), exp: ""},
		{impl: dbutil.Postgres, time: time.Unix(0, 1620721781789035200), exp: ""},
		{impl: dbutil.Cockroach, time: time.Unix(0, 1620721781789035200), exp: " AS OF SYSTEM TIME '1620721781789035200' "},
		{impl: dbutil.Bolt, time: time.Unix(0, 1620721781789035200), exp: ""},
		{impl: dbutil.Redis, time: time.Unix(0, 1620721781789035200), exp: ""},
		{impl: dbutil.SQLite3, time: time.Unix(0, 1620721781789035200), exp: ""},
	}

	for _, test := range tests {
		asof := test.impl.AsOfSystemTime(test.time)
		assert.Equal(t, test.exp, asof)
	}
}

func TestAsOfSystemInterval(t *testing.T) {
	tests := []struct {
		impl     dbutil.Implementation
		interval time.Duration
		exp      string
	}{
		{impl: dbutil.Unknown, interval: 0, exp: ""},
		{impl: dbutil.Postgres, interval: 0, exp: ""},
		{impl: dbutil.Cockroach, interval: 0, exp: ""},
		{impl: dbutil.Bolt, interval: 0, exp: ""},
		{impl: dbutil.Redis, interval: 0, exp: ""},
		{impl: dbutil.SQLite3, interval: 0, exp: ""},

		{impl: dbutil.Unknown, interval: 1, exp: ""},
		{impl: dbutil.Postgres, interval: 1, exp: ""},
		{impl: dbutil.Cockroach, interval: 1, exp: ""},
		{impl: dbutil.Bolt, interval: 1, exp: ""},
		{impl: dbutil.Redis, interval: 1, exp: ""},
		{impl: dbutil.SQLite3, interval: 1, exp: ""},

		{impl: dbutil.Unknown, interval: -1, exp: ""},
		{impl: dbutil.Postgres, interval: -1, exp: ""},
		{impl: dbutil.Cockroach, interval: -time.Nanosecond, exp: " AS OF SYSTEM TIME '-1µs' "},
		{impl: dbutil.Bolt, interval: -1, exp: ""},
		{impl: dbutil.Redis, interval: -1, exp: ""},
		{impl: dbutil.SQLite3, interval: -1, exp: ""},

		{impl: dbutil.Unknown, interval: -time.Millisecond, exp: ""},
		{impl: dbutil.Postgres, interval: -time.Millisecond, exp: ""},
		{impl: dbutil.Cockroach, interval: -time.Millisecond, exp: " AS OF SYSTEM TIME '-1ms' "},
		{impl: dbutil.Bolt, interval: -time.Millisecond, exp: ""},
		{impl: dbutil.Redis, interval: -time.Millisecond, exp: ""},
		{impl: dbutil.SQLite3, interval: -time.Millisecond, exp: ""},
	}

	for _, test := range tests {
		asof := test.impl.AsOfSystemInterval(test.interval)
		assert.Equal(t, test.exp, asof)
	}
}
