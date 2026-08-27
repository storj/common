// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import "storj.io/picobuf"

// This file implements picobuf field encoding for the types that are used as
// (pico.field).custom_type in the protobuf definitions, the same role the
// Marshal/MarshalTo/Unmarshal/Size methods play for gogo's customtype.
//
// The encoding matches gogo/protobuf field for field, because order limits,
// piece hashes and stream and segment IDs are signed by serializing them and a
// peer of a different version has to arrive at the same bytes. gogo emits a
// non-nullable custom type unconditionally, including when it is the zero
// value, except when the Go type is a pointer, slice or map and the value is
// nil, which it skips (proto.marshalInfo.marshal: "nil pointer always marshals
// to nothing"). So the fixed size types and the piece keys always encode, and
// StreamID and SegmentID encode only when non-nil, an empty non-nil one
// included. Every encoder also tolerates a nil receiver, so a field declared
// as a pointer is skipped rather than panicking.
//
// One case cannot be reproduced from here: gogo decides per field declaration,
// so a StreamID field declared as a pointer that points at a nil slice encodes
// as a zero length field there, while it encodes as nothing here. Only
// ObjectListItem.stream_id and ListPendingObjectStreamsResponse are declared
// that way, no signed message contains such a field, and a zero length field
// and an absent one decode to the same value.
//
// The decoders copy out of the input rather than retaining the slice the
// decoder hands them: these values outlive the RPC they arrived on, and the
// buffer underneath is reused for the next packet on the connection. BytesView
// borrows each field temporarily; Unmarshal makes the one required copy into
// the destination, without an intermediate byte-slice allocation.

var (
	_ picobuf.CustomType = (*NodeID)(nil)
	_ picobuf.CustomType = (*PieceID)(nil)
	_ picobuf.CustomType = (*Nonce)(nil)
	_ picobuf.CustomType = (*SerialNumber)(nil)
	_ picobuf.CustomType = (*StreamID)(nil)
	_ picobuf.CustomType = (*SegmentID)(nil)
	_ picobuf.CustomType = (*PiecePublicKey)(nil)
	_ picobuf.CustomType = (*PiecePrivateKey)(nil)
)

// PicoEncode implements picobuf custom encoding.
func (id *NodeID) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if id == nil {
		return
	}
	v := id[:]
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (id *NodeID) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := id.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (id *PieceID) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if id == nil {
		return
	}
	v := id[:]
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (id *PieceID) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := id.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (nonce *Nonce) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if nonce == nil {
		return
	}
	v := nonce[:]
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (nonce *Nonce) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := nonce.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (id *SerialNumber) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if id == nil {
		return
	}
	v := id[:]
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (id *SerialNumber) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := id.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (id *StreamID) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if id == nil || *id == nil {
		return
	}
	v := []byte(*id)
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (id *StreamID) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := id.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (id *SegmentID) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if id == nil || *id == nil {
		return
	}
	v := []byte(*id)
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (id *SegmentID) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := id.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (key *PiecePublicKey) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if key == nil {
		return
	}
	v := []byte(key.pub)
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (key *PiecePublicKey) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := key.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}

// PicoEncode implements picobuf custom encoding.
func (key *PiecePrivateKey) PicoEncode(c *picobuf.Encoder, field picobuf.FieldNumber) {
	if key == nil {
		return
	}
	v := []byte(key.priv)
	c.AlwaysBytes(field, &v)
}

// PicoDecode implements picobuf custom decoding.
func (key *PiecePrivateKey) PicoDecode(c *picobuf.Decoder, field picobuf.FieldNumber) {
	if !c.MatchField(field, 2) {
		return
	}
	var v []byte
	c.BytesView(field, &v)
	if err := key.Unmarshal(v); err != nil {
		c.Fail(field, err.Error())
	}
}
