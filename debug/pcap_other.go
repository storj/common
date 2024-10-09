// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build !linux

package debug

import (
	"context"
	"sync/atomic"
)

func capturePackets(ctx context.Context, stop *atomic.Bool) {}
