// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// Package pbgrpc contains grpc definitions for Storj Network.
package pbgrpc

import (
	context "context"
	"crypto/tls"
	"net"

	"github.com/zeebo/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"storj.io/common/internal/grpchook"
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
}
