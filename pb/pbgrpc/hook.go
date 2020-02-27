// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package pbgrpc contains grpc definitions for Storj Network.
package pbgrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/zeebo/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"storj.io/common/internal/grpchook"
	"storj.io/common/rpc/rpcstatus"
)

// Error is the default error class for grpchook.
var Error = errs.Class("grpchook")

// Here we hook grpc into rest of the systems.
func init() {
	grpchook.HookedErrServerStopped = grpc.ErrServerStopped

	// InternalFromContext returns the peer that was previously associated by NewContext using grpc.
	grpchook.HookedInternalFromContext = func(ctx context.Context) (addr net.Addr, state tls.ConnectionState, err error) {
		peer, ok := peer.FromContext(ctx)
		if !ok {
			return nil, tls.ConnectionState{}, Error.New("unable to get grpc peer from context")
		}

		tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
		if !ok {
			return nil, tls.ConnectionState{}, Error.New("peer AuthInfo is not credentials.TLSInfo")
		}

		return peer.Addr, tlsInfo.State, nil
	}

	grpchook.HookedConvertToStatusCode = func(err error) (grpchook.StatusCode, bool) {
		code := status.Code(err)
		if code == codes.Unknown {
			return grpchook.StatusCode(rpcstatus.Unknown), false
		}

		convertedCode := statusCodeFromGRPC(code)
		return grpchook.StatusCode(convertedCode), true
	}

	grpchook.HookedErrorWrap = func(hookcode grpchook.StatusCode, err error) error {
		if err == nil {
			return nil
		}
		code := rpcstatus.StatusCode(hookcode)

		ce := &codeErr{
			code: code,
			grpc: status.New(statusCodeToGRPC(code), err.Error()),
		}

		if ee, ok := err.(errsError); ok {
			ce.errsError = ee
		} else {
			ce.errsError = errs.Wrap(err).(errsError)
		}

		return ce
	}
}

type errsError interface {
	error
	fmt.Formatter
	Name() (string, bool)
}

// codeErr implements error that can work both in grpc and drpc.
type codeErr struct {
	errsError
	code rpcstatus.StatusCode
	grpc *status.Status
}

func (c *codeErr) Unwrap() error { return c.errsError }
func (c *codeErr) Cause() error  { return c.errsError }

func (c *codeErr) Code() uint64               { return uint64(c.code) }
func (c *codeErr) GRPCStatus() *status.Status { return c.grpc }

// statusCodeToGRPC returns the grpc version of the status code.
func statusCodeToGRPC(s rpcstatus.StatusCode) codes.Code {
	switch s {
	case rpcstatus.Unknown:
		return codes.Unknown
	case rpcstatus.OK:
		return codes.OK
	case rpcstatus.Canceled:
		return codes.Canceled
	case rpcstatus.InvalidArgument:
		return codes.InvalidArgument
	case rpcstatus.DeadlineExceeded:
		return codes.DeadlineExceeded
	case rpcstatus.NotFound:
		return codes.NotFound
	case rpcstatus.AlreadyExists:
		return codes.AlreadyExists
	case rpcstatus.PermissionDenied:
		return codes.PermissionDenied
	case rpcstatus.ResourceExhausted:
		return codes.ResourceExhausted
	case rpcstatus.FailedPrecondition:
		return codes.FailedPrecondition
	case rpcstatus.Aborted:
		return codes.Aborted
	case rpcstatus.OutOfRange:
		return codes.OutOfRange
	case rpcstatus.Unimplemented:
		return codes.Unimplemented
	case rpcstatus.Internal:
		return codes.Internal
	case rpcstatus.Unavailable:
		return codes.Unavailable
	case rpcstatus.DataLoss:
		return codes.DataLoss
	case rpcstatus.Unauthenticated:
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}

// statusCodeFromGRPC turns the grpc status code into a StatusCode.
func statusCodeFromGRPC(code codes.Code) rpcstatus.StatusCode {
	switch code {
	case codes.Unknown:
		return rpcstatus.Unknown
	case codes.OK:
		return rpcstatus.OK
	case codes.Canceled:
		return rpcstatus.Canceled
	case codes.InvalidArgument:
		return rpcstatus.InvalidArgument
	case codes.DeadlineExceeded:
		return rpcstatus.DeadlineExceeded
	case codes.NotFound:
		return rpcstatus.NotFound
	case codes.AlreadyExists:
		return rpcstatus.AlreadyExists
	case codes.PermissionDenied:
		return rpcstatus.PermissionDenied
	case codes.ResourceExhausted:
		return rpcstatus.ResourceExhausted
	case codes.FailedPrecondition:
		return rpcstatus.FailedPrecondition
	case codes.Aborted:
		return rpcstatus.Aborted
	case codes.OutOfRange:
		return rpcstatus.OutOfRange
	case codes.Unimplemented:
		return rpcstatus.Unimplemented
	case codes.Internal:
		return rpcstatus.Internal
	case codes.Unavailable:
		return rpcstatus.Unavailable
	case codes.DataLoss:
		return rpcstatus.DataLoss
	case codes.Unauthenticated:
		return rpcstatus.Unauthenticated
	default:
		return rpcstatus.Unknown
	}
}
