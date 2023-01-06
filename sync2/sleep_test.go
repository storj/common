// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"context"
	"testing"
	"time"

	"storj.io/common/sync2"
	"storj.io/common/time2"
)

func TestSleep(t *testing.T) {
	t.Parallel()

	t.Run("against the real clock", func(t *testing.T) {
		const sleepError = time.Second / 2 // should be larger than most system error with regards to sleep

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		start := time.Now()
		if !sync2.Sleep(ctx, time.Second) {
			t.Error("expected true as result")
		}
		if time.Since(start) < time.Second-sleepError {
			t.Error("sleep took too little time")
		}
	})

	t.Run("against a fake clock", func(t *testing.T) {
		ctx, timeMachine := time2.WithNewMachine(context.Background())

		defer sync2.Go(func() { timeMachine.BlockThenAdvance(ctx, 1, time.Second) })()

		start := timeMachine.Now()
		if !sync2.Sleep(ctx, time.Second) {
			t.Error("expected true as result")
		}
		if timeMachine.Since(start) != time.Second {
			t.Error("sleep took too little time")
		}
	})
}

func TestSleep_Cancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	if sync2.Sleep(ctx, 5*time.Second) {
		t.Error("expected false as result")
	}
	if time.Since(start) > time.Second {
		t.Error("sleep took too long")
	}
}
