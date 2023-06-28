package nodetag

import (
	"bytes"
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/zeebo/errs"

	"storj.io/common/pb"
	"storj.io/common/signing"
)

var (
	SignatureErr     = errs.Class("invalid signature")
	SerializationErr = errs.Class("invalid tag serialization")
	WrongSignee      = errs.Class("node id mismatch")
)

// Sign create a signed tag set from a raw one.
func Sign(ctx context.Context, tagSet *pb.NodeTagSet, signer signing.Signer) (*pb.SignedNodeTagSet, error) {
	signed := &pb.SignedNodeTagSet{}
	raw, err := proto.Marshal(tagSet)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	signature, err := signer.HashAndSign(ctx, raw)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	signed.Signature = signature
	signed.SignerNodeId = signer.ID().Bytes()
	signed.SerializedTag = raw
	return signed, nil
}

// Verify checks the signature of a signed tag set.
func Verify(ctx context.Context, tags *pb.SignedNodeTagSet, signee signing.Signee) (*pb.NodeTagSet, error) {
	if !bytes.Equal(tags.SignerNodeId, signee.ID().Bytes()) {
		return nil, WrongSignee.New("wrong signee to verify")
	}
	err := signee.HashAndVerifySignature(ctx, tags.SerializedTag, tags.Signature)
	if err != nil {
		return nil, SignatureErr.Wrap(err)
	}
	tagset := &pb.NodeTagSet{}
	err = proto.Unmarshal(tags.SerializedTag, tagset)
	if err != nil {
		return nil, SerializationErr.Wrap(err)
	}
	return tagset, nil
}
