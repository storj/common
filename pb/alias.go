// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package pb

import "storj.io/picobuf"

// Message is an alias for picobuf.Message.
type Message = picobuf.Message

// Unmarshal is an alias for picobuf.Unmarshal.
func Unmarshal(buf []byte, pb Message) error { return picobuf.Unmarshal(buf, pb) }

// Marshal encodes pb with map entries sorted by key.
//
// Signatures are verified over re-encoded messages, so the encoding must not
// depend on map iteration order.
func Marshal(pb Message) ([]byte, error) {
	return picobuf.MarshalOptions{Deterministic: true}.Marshal(pb)
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
