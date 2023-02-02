// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package pb

import (
	"context"

	"storj.io/drpc"
)

// DRPCReplaySafePiecestoreClient is a client that exposes the replay safe
// subset of DRPCPiecestoreClient methods.
type DRPCReplaySafePiecestoreClient struct {
	c DRPCPiecestoreClient
}

// NewDRPCReplaySafePiecestoreClient makes a DRPCReplaySafePiecestoreClient.
func NewDRPCReplaySafePiecestoreClient(cc drpc.Conn) *DRPCReplaySafePiecestoreClient {
	return &DRPCReplaySafePiecestoreClient{c: NewDRPCPiecestoreClient(cc)}
}

// DRPCConn returns the drpc.Conn.
func (c *DRPCReplaySafePiecestoreClient) DRPCConn() drpc.Conn { return c.c.DRPCConn() }

// Upload calls DRPCPiecestoreClient.Upload.
func (c *DRPCReplaySafePiecestoreClient) Upload(ctx context.Context) (DRPCPiecestore_UploadClient, error) {
	return c.c.Upload(ctx)
}

// Download calls DRPCPiecestoreClient.Download.
func (c *DRPCReplaySafePiecestoreClient) Download(ctx context.Context) (DRPCPiecestore_DownloadClient, error) {
	return c.c.Download(ctx)
}

// DRPCReplaySafePiecestoreServer is a server that exposes the replay safe
// subset of DRPCPiecestoreServer methods.
type DRPCReplaySafePiecestoreServer interface {
	Upload(DRPCPiecestore_UploadStream) error
	Download(DRPCPiecestore_DownloadStream) error
}

type drpcReplaySafePiecestoreDescription struct{}

func (drpcReplaySafePiecestoreDescription) NumMethods() int { return 2 }

func (drpcReplaySafePiecestoreDescription) Method(n int) (string, drpc.Encoding, drpc.Receiver, interface{}, bool) {
	switch n {
	case 0:
		rpc, enc, receiver, method, ok := (DRPCPiecestoreDescription{}).Method(0)
		return rpc, enc, receiver, method, ok && rpc == "/piecestore.Piecestore/Upload"
	case 1:
		rpc, enc, receiver, method, ok := (DRPCPiecestoreDescription{}).Method(1)
		return rpc, enc, receiver, method, ok && rpc == "/piecestore.Piecestore/Download"
	default:
		return "", nil, nil, nil, false
	}
}

// DRPCRegisterReplaySafePiecestore registers a replay safe Piecestore Server on
// the provided drpc.Mux.
func DRPCRegisterReplaySafePiecestore(mux drpc.Mux, impl DRPCReplaySafePiecestoreServer) error {
	return mux.Register(impl, drpcReplaySafePiecestoreDescription{})
}
