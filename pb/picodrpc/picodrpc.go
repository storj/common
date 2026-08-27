// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

// Package picodrpc adapts picobuf to the encoding interface that
// protoc-gen-go-drpc generates against, so that the generated code can use
// picobuf without drpc needing to know about it.
package picodrpc

import (
	"storj.io/drpc"
	"storj.io/picobuf"
)

// Marshal returns the encoded form of msg.
func Marshal(msg drpc.Message) ([]byte, error) {
	return picobuf.Marshal(msg.(picobuf.Message))
}

// MarshalAppend returns the encoded form of msg appended to buf.
func MarshalAppend(buf []byte, msg drpc.Message) ([]byte, error) {
	return picobuf.MarshalAppend(buf, msg.(picobuf.Message))
}

// Unmarshal reads the encoded form of a single message into msg.
func Unmarshal(buf []byte, msg drpc.Message) error {
	return picobuf.Unmarshal(buf, msg.(picobuf.Message))
}
