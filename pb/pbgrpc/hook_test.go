// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package pbgrpc contains grpc definitions for Storj Network.
package pbgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	"storj.io/common/rpc/rpcstatus"
	"storj.io/drpc/drpcerr"
)

var allCodes = []rpcstatus.StatusCode{
	rpcstatus.Unknown,
	rpcstatus.OK,
	rpcstatus.Canceled,
	rpcstatus.InvalidArgument,
	rpcstatus.DeadlineExceeded,
	rpcstatus.NotFound,
	rpcstatus.AlreadyExists,
	rpcstatus.PermissionDenied,
	rpcstatus.ResourceExhausted,
	rpcstatus.FailedPrecondition,
	rpcstatus.Aborted,
	rpcstatus.OutOfRange,
	rpcstatus.Unimplemented,
	rpcstatus.Internal,
	rpcstatus.Unavailable,
	rpcstatus.DataLoss,
	rpcstatus.Unauthenticated,
}

func TestStatus(t *testing.T) {
	for _, code := range allCodes {
		err := rpcstatus.Error(code, "")
		assert.Equal(t, rpcstatus.Code(err), code)
		assert.Equal(t, status.Code(err), statusCodeToGRPC(code))
		assert.Equal(t, drpcerr.Code(err), uint64(code))
	}

	assert.Equal(t, rpcstatus.Code(nil), rpcstatus.OK)
	assert.Equal(t, rpcstatus.Code(context.Canceled), rpcstatus.Canceled)
	assert.Equal(t, rpcstatus.Code(context.DeadlineExceeded), rpcstatus.DeadlineExceeded)
}
