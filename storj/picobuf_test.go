// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/picobuf"
)

const customField = 3

type picoField[T any] interface {
	*T
	PicoEncode(*picobuf.Encoder, picobuf.FieldNumber)
	PicoDecode(*picobuf.Decoder, picobuf.FieldNumber)
}

// customMessage is a message with a single custom typed field.
type customMessage[T any, P picoField[T]] struct{ Value T }

func (m *customMessage[T, P]) Encode(c *picobuf.Encoder) bool {
	P(&m.Value).PicoEncode(c, customField)
	return true
}

func (m *customMessage[T, P]) Decode(c *picobuf.Decoder) {
	P(&m.Value).PicoDecode(c, customField)
}

// checkCustom verifies that value encodes as wantPayload in a length delimited
// field, which is emitted even when the value is zero, and that it survives a
// round-trip.
func checkCustom[T any, P picoField[T]](t *testing.T, value T, wantPayload []byte) {
	t.Helper()

	data, err := picobuf.Marshal(&customMessage[T, P]{Value: value})
	require.NoError(t, err)

	want := append([]byte{customField<<3 | 2, byte(len(wantPayload))}, wantPayload...)
	require.Equal(t, want, data)

	for _, alias := range []bool{false, true} {
		input := append([]byte(nil), data...)
		var out customMessage[T, P]
		require.NoError(t, (picobuf.UnmarshalOptions{AliasInput: alias}).Unmarshal(input, &out))
		require.Equal(t, value, out.Value)
		clear(input)
		require.Equal(t, value, out.Value, "custom decoder retained input (AliasInput=%v)", alias)
	}
}

func TestCustomTypes(t *testing.T) {
	t.Run("NodeID", func(t *testing.T) {
		checkCustom[NodeID](t, NodeID{}, make([]byte, NodeIDSize))
		id := NodeID{1, 2, 3, 31: 8}
		checkCustom[NodeID](t, id, id[:])
	})

	t.Run("PieceID", func(t *testing.T) {
		checkCustom[PieceID](t, PieceID{}, make([]byte, len(PieceID{})))
		id := NewPieceID()
		checkCustom[PieceID](t, id, id[:])
	})

	t.Run("Nonce", func(t *testing.T) {
		checkCustom[Nonce](t, Nonce{}, make([]byte, NonceSize))
		nonce := Nonce{1, 2, 3}
		checkCustom[Nonce](t, nonce, nonce[:])
	})

	t.Run("SerialNumber", func(t *testing.T) {
		checkCustom[SerialNumber](t, SerialNumber{}, make([]byte, len(SerialNumber{})))
		serial := SerialNumber{1, 2, 3}
		checkCustom[SerialNumber](t, serial, serial[:])
	})

	// a nil StreamID or SegmentID is skipped, matching gogo, which skips a
	// nil slice even for a non-nullable custom type. An empty non-nil one is
	// still encoded, as a zero length field.
	t.Run("StreamID", func(t *testing.T) {
		checkAbsent[StreamID](t, nil)
		checkCustom[StreamID](t, StreamID{}, []byte{})
		checkCustom[StreamID](t, StreamID{1, 2, 3}, []byte{1, 2, 3})
	})

	t.Run("SegmentID", func(t *testing.T) {
		checkAbsent[SegmentID](t, nil)
		checkCustom[SegmentID](t, SegmentID{}, []byte{})
		checkCustom[SegmentID](t, SegmentID{1, 2, 3}, []byte{1, 2, 3})
	})

	t.Run("PieceKey", func(t *testing.T) {
		checkCustom[PiecePublicKey](t, PiecePublicKey{}, []byte{})
		checkCustom[PiecePrivateKey](t, PiecePrivateKey{}, []byte{})

		pub, priv, err := NewPieceKey()
		require.NoError(t, err)
		checkCustom[PiecePublicKey](t, pub, pub.Bytes())
		checkCustom[PiecePrivateKey](t, priv, priv.Bytes())
	})
}

// checkAbsent verifies that value is not encoded at all.
func checkAbsent[T any, P picoField[T]](t *testing.T, value T) {
	t.Helper()

	data, err := picobuf.Marshal(&customMessage[T, P]{Value: value})
	require.NoError(t, err)
	require.Empty(t, data)
}

// TestCustomTypesNilReceiver verifies that the encoders tolerate a nil
// receiver, which is how a field declared as a pointer is skipped.
func TestCustomTypesNilReceiver(t *testing.T) {
	encoders := []func(*picobuf.Encoder, picobuf.FieldNumber){
		(*NodeID)(nil).PicoEncode,
		(*PieceID)(nil).PicoEncode,
		(*Nonce)(nil).PicoEncode,
		(*SerialNumber)(nil).PicoEncode,
		(*StreamID)(nil).PicoEncode,
		(*SegmentID)(nil).PicoEncode,
		(*PiecePublicKey)(nil).PicoEncode,
		(*PiecePrivateKey)(nil).PicoEncode,
	}
	for _, encode := range encoders {
		// a nil encoder panics unless the receiver check returns first.
		encode(nil, 1)
	}
}

// TestCustomTypesDecodeCopies verifies that decoding does not retain the input,
// which is reused for the next packet on a drpc connection.
func TestCustomTypesDecodeCopies(t *testing.T) {
	data := []byte{customField<<3 | 2, 3, 1, 2, 3}

	var stream customMessage[StreamID, *StreamID]
	require.NoError(t, picobuf.UnmarshalOptions{AliasInput: true}.Unmarshal(data, &stream))
	var segment customMessage[SegmentID, *SegmentID]
	require.NoError(t, picobuf.UnmarshalOptions{AliasInput: true}.Unmarshal(data, &segment))

	clear(data)
	require.Equal(t, StreamID{1, 2, 3}, stream.Value)
	require.Equal(t, SegmentID{1, 2, 3}, segment.Value)
}

func TestCustomTypesInvalid(t *testing.T) {
	data := []byte{customField<<3 | 2, 3, 1, 2, 3}

	var node customMessage[NodeID, *NodeID]
	require.Error(t, picobuf.Unmarshal(data, &node))

	var piece customMessage[PieceID, *PieceID]
	require.Error(t, picobuf.Unmarshal(data, &piece))

	var nonce customMessage[Nonce, *Nonce]
	require.Error(t, picobuf.Unmarshal(data, &nonce))

	var serial customMessage[SerialNumber, *SerialNumber]
	require.Error(t, picobuf.Unmarshal(data, &serial))

	var pub customMessage[PiecePublicKey, *PiecePublicKey]
	require.Error(t, picobuf.Unmarshal(data, &pub))

	var priv customMessage[PiecePrivateKey, *PiecePrivateKey]
	require.Error(t, picobuf.Unmarshal(data, &priv))
}

// TestCustomTypesAbsent verifies that a field that is not in the input leaves
// the value untouched rather than failing to decode.
func TestCustomTypesAbsent(t *testing.T) {
	var out customMessage[NodeID, *NodeID]
	require.NoError(t, picobuf.Unmarshal(nil, &out))
	require.Equal(t, NodeID{}, out.Value)

	other := []byte{1<<3 | 2, 1, 9}
	require.NoError(t, picobuf.Unmarshal(other, &out))
	require.Equal(t, NodeID{}, out.Value)
}

// TestCustomTypesWrongWireType verifies that a field with the right number but
// a wire type other than bytes is skipped and leaves the value untouched.
func TestCustomTypesWrongWireType(t *testing.T) {
	pub, priv, err := NewPieceKey()
	require.NoError(t, err)

	checkWrongWireType(t, NodeID{1})
	checkWrongWireType(t, PieceID{1})
	checkWrongWireType(t, Nonce{1})
	checkWrongWireType(t, SerialNumber{1})
	checkWrongWireType(t, StreamID{1})
	checkWrongWireType(t, SegmentID{1})
	checkWrongWireType(t, pub)
	checkWrongWireType(t, priv)
}

func checkWrongWireType[T any, P picoField[T]](t *testing.T, value T) {
	t.Helper()

	varint := []byte{customField << 3, 5} // wire type 0, a varint
	out := customMessage[T, P]{Value: value}
	require.NoError(t, picobuf.Unmarshal(varint, &out))
	require.Equal(t, value, out.Value)
}
