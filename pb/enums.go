// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package pb

// The generated code does not include the name and value maps that
// gogo/protobuf used to emit. These are kept for the enumerations that
// dependents still iterate over or look up. Keep them in sync with orders.proto.

// PieceAction_name maps a PieceAction to its name.
var PieceAction_name = map[int32]string{
	0: "INVALID",
	1: "PUT",
	2: "GET",
	3: "GET_AUDIT",
	4: "GET_REPAIR",
	5: "PUT_REPAIR",
	6: "DELETE",
	7: "PUT_GRACEFUL_EXIT",
}

// PieceAction_value maps a PieceAction name to its value.
var PieceAction_value = map[string]int32{
	"INVALID":           0,
	"PUT":               1,
	"GET":               2,
	"GET_AUDIT":         3,
	"GET_REPAIR":        4,
	"PUT_REPAIR":        5,
	"DELETE":            6,
	"PUT_GRACEFUL_EXIT": 7,
}

// PieceHashAlgorithm_name maps a PieceHashAlgorithm to its name.
var PieceHashAlgorithm_name = map[int32]string{
	0: "SHA256",
	1: "BLAKE3",
}

// PieceHashAlgorithm_value maps a PieceHashAlgorithm name to its value.
var PieceHashAlgorithm_value = map[string]int32{
	"SHA256": 0,
	"BLAKE3": 1,
}

// SettlementWithWindowResponse_Status_name maps a SettlementWithWindowResponse_Status to its name.
var SettlementWithWindowResponse_Status_name = map[int32]string{
	0: "ACCEPTED",
	1: "REJECTED",
}

// SettlementWithWindowResponse_Status_value maps a SettlementWithWindowResponse_Status name to its value.
var SettlementWithWindowResponse_Status_value = map[string]int32{
	"ACCEPTED": 0,
	"REJECTED": 1,
}
