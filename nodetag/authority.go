package nodetag

import (
	"bytes"
	"context"
	"github.com/zeebo/errs"
	"storj.io/common/pb"
	"storj.io/common/signing"
)

var (
	UnkownSignee = errs.Class("node tag signee is unknown")
)

// Authority contains all possible signee.
type Authority []signing.Signee

// Verify checks if any of the storage signee can validate the signature.
func (a Authority) Verify(ctx context.Context, tags *pb.SignedNodeTagSet) (*pb.NodeTagSet, error) {
	for _, signee := range a {
		if bytes.Equal(signee.ID().Bytes(), tags.SignerNodeId) {
			return Verify(ctx, tags, signee)
		}
	}
	return nil, UnkownSignee.New("no certificate for signer nodeID: %x", tags.SignerNodeId)
}
