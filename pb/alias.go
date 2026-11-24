// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package pb

import "storj.io/picobuf"

// Message is an alias for picobuf.Message.
type Message = picobuf.Message

// Unmarshal resets pb and decodes buf into it.
//
// picobuf.Unmarshal merges into pb, gogo/protobuf reset it first; callers
// reuse messages and rely on the reset.
func Unmarshal[T any, PT interface {
	*T
	Message
}](buf []byte, pb PT) error {
	var zero T
	*pb = zero
	return picobuf.Unmarshal(buf, pb)
}

// Marshal encodes pb with map entries sorted by key.
//
// Signatures are verified over re-encoded messages, so the encoding must not
// depend on map iteration order.
//
// An empty message encodes to a non-nil empty slice, as with gogo/protobuf,
// so that it is not stored as NULL in a database.
func Marshal(pb Message) ([]byte, error) {
	data, err := picobuf.MarshalOptions{Deterministic: true}.Marshal(pb)
	if data == nil && err == nil {
		data = []byte{}
	}
	return data, err
}

// Size returns the encoded size of pb.
//
// It returns zero when pb cannot be encoded, which is indistinguishable from
// an empty message, so use it only where an approximate size is good enough,
// such as telemetry. Use Marshal when the error matters.
func Size(pb Message) int {
	size, err := picobuf.Size(pb)
	if err != nil {
		return 0
	}
	return size
}
