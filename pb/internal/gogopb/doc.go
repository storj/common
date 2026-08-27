// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

// Package gogopb contains the gogo/protobuf generated code for the same
// protobuf definitions as storj.io/common/pb, from before the migration to
// picobuf. It exists only as a test oracle: the tests in pb compare the
// picobuf encoding against it byte for byte, so that a mistake in the
// annotations or a custom type cannot silently change the wire format.
//
// The .proto files here are copies that keep the original gogoproto
// annotations, so that a mistake in the pico annotations shows up as an
// encoding difference. When a definition in pb changes, copy the change here
// with the equivalent gogoproto annotations and run go generate in pb, which
// regenerates this package. TestGogoProtoInSync fails until the two match. Delete
// this package once the migration has been in production for a release.
package gogopb
