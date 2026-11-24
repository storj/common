// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package pb_test

import (
	"bytes"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	gogo "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"storj.io/common/pb"
	"storj.io/common/storj"
	"storj.io/picobuf"
)

// gogoMessage is the gogo/protobuf side of a message pair.
type gogoMessage interface {
	gogo.Message
}

type messagePair struct {
	name string
	pico func() pb.Message
	gogo func() gogoMessage
}

// TestGogoCompatibility fills the picobuf and the gogo/protobuf version of
// every message with the same pseudo-random values and requires that they
// encode to the same bytes and can decode each other's output. This is what
// keeps the migration from picobuf silently changing the wire format, which
// would break signature verification.
func TestGogoCompatibility(t *testing.T) {
	for _, message := range messages {
		t.Run(message.name, func(t *testing.T) {
			// an untouched message first: this is what covers the always
			// present fields at their zero value, which are encoded rather
			// than omitted.
			requireSameEncoding(t, message.pico(), message.gogo(), message.pico())

			for seed := int64(0); seed < 20; seed++ {
				picoMsg := message.pico()
				gogoMsg := message.gogo()

				fill(t, rand.New(rand.NewSource(seed)), reflect.ValueOf(picoMsg).Elem(), 0)
				fill(t, rand.New(rand.NewSource(seed)), reflect.ValueOf(gogoMsg).Elem(), 0)

				requireSameEncoding(t, picoMsg, gogoMsg, message.pico())
			}
		})
	}
}

// requireSameEncoding requires that the picobuf and the gogo/protobuf message
// encode to the same bytes, and that each side can decode and re-encode what
// the other wrote. empty is an unfilled message of the picobuf type.
func requireSameEncoding(t *testing.T, picoMsg pb.Message, gogoMsg gogoMessage, empty pb.Message) {
	t.Helper()

	picoData, err := pb.Marshal(picoMsg)
	require.NoError(t, err)
	gogoData, err := gogo.Marshal(gogoMsg)
	require.NoError(t, err)
	require.True(t, bytes.Equal(gogoData, picoData), "encoding differs:\n%x\n%x", gogoData, picoData)

	// picobuf must decode what gogo wrote and re-encode it unchanged, which is
	// what signature verification across version skew relies on.
	require.NoError(t, picobuf.Unmarshal(gogoData, empty))
	redone, err := pb.Marshal(empty)
	require.NoError(t, err)
	require.True(t, bytes.Equal(gogoData, redone), "picobuf round-trip differs:\n%x\n%x", gogoData, redone)

	// and gogo must decode what picobuf wrote.
	decodedGogo := reflect.New(reflect.TypeOf(gogoMsg).Elem()).Interface().(gogoMessage)
	require.NoError(t, gogo.Unmarshal(picoData, decodedGogo))
	redone, err = gogo.Marshal(decodedGogo)
	require.NoError(t, err)
	require.True(t, bytes.Equal(picoData, redone), "gogo round-trip differs:\n%x\n%x", picoData, redone)
}

var (
	typeTime            = reflect.TypeOf(time.Time{})
	typeDuration        = reflect.TypeOf(time.Duration(0))
	typePiecePublicKey  = reflect.TypeOf(storj.PiecePublicKey{})
	typePiecePrivateKey = reflect.TypeOf(storj.PiecePrivateKey{})
)

// maxDepth bounds the recursive messages, e.g. metainfo batch requests.
const maxDepth = 3

// fill assigns pseudo-random values to v. Two structurally identical types
// filled from the same seed get the same values, which is what lets the picobuf
// and the gogo/protobuf message be compared.
func fill(t *testing.T, rng *rand.Rand, v reflect.Value, depth int) {
	t.Helper()

	switch v.Type() {
	case typeTime:
		// a zero time must be covered: an always_present timestamp encodes it
		// as seconds=-62135596800 rather than omitting the field, and a
		// timestamp that is absent in practice, such as an order limit without
		// a piece expiration, is exactly that case.
		switch rng.Intn(4) {
		case 0:
			v.Set(reflect.ValueOf(time.Time{}))
		case 1:
			v.Set(reflect.ValueOf(time.Unix(-rng.Int63n(1<<32), rng.Int63n(1e9)).UTC()))
		default:
			// keep well inside the range that google.protobuf.Timestamp allows.
			v.Set(reflect.ValueOf(time.Unix(rng.Int63n(1<<32), rng.Int63n(1e9)).UTC()))
		}
		return
	case typeDuration:
		v.Set(reflect.ValueOf(time.Duration(rng.Int63n(1e12))))
		return
	case typePiecePublicKey:
		key, err := storj.PiecePublicKeyFromBytes(randBytes(rng, 32))
		require.NoError(t, err)
		v.Set(reflect.ValueOf(key))
		return
	case typePiecePrivateKey:
		key, err := storj.PiecePrivateKeyFromBytes(randBytes(rng, 64))
		require.NoError(t, err)
		v.Set(reflect.ValueOf(key))
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(rng.Intn(2) == 1)
	case reflect.Int32, reflect.Int64, reflect.Int:
		// enums are named int32 types; keep the values small so that they stay
		// inside the declared range.
		if v.Type().Name() != "int32" && v.Kind() == reflect.Int32 {
			v.SetInt(rng.Int63n(4))
		} else {
			v.SetInt(rng.Int63n(1 << 40))
		}
	case reflect.Uint32, reflect.Uint64, reflect.Uint:
		v.SetUint(uint64(rng.Int63n(1 << 40)))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(rng.Int31n(1 << 20)))
	case reflect.String:
		v.SetString(randString(rng))
	case reflect.Array:
		// fixed size byte arrays, e.g. storj.NodeID.
		require.Equal(t, reflect.Uint8, v.Type().Elem().Kind(), "unhandled array %v", v.Type())
		for i := 0; i < v.Len(); i++ {
			v.Index(i).SetUint(uint64(rng.Intn(256)))
		}
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// nil, empty and non-empty are three different cases for a custom
			// typed field such as storj.StreamID: gogo skips a nil slice and
			// encodes an empty non-nil one as a zero length field.
			switch rng.Intn(6) {
			case 0:
				return
			case 1:
				v.SetBytes([]byte{})
			default:
				v.SetBytes(randBytes(rng, 1+rng.Intn(16)))
			}
			return
		}
		if depth >= maxDepth {
			return
		}
		n := 1 + rng.Intn(2)
		v.Set(reflect.MakeSlice(v.Type(), n, n))
		for i := 0; i < n; i++ {
			elem := v.Index(i)
			// a repeated message field may not contain nil elements.
			if elem.Kind() == reflect.Pointer {
				elem.Set(reflect.New(elem.Type().Elem()))
				fill(t, rng, elem.Elem(), depth+1)
				continue
			}
			fill(t, rng, elem, depth)
		}
	case reflect.Map:
		// a single entry only: gogo/protobuf encodes maps in Go's iteration
		// order while picobuf sorts the keys, so more than one entry is not
		// comparable byte for byte.
		v.Set(reflect.MakeMap(v.Type()))
		key := reflect.New(v.Type().Key()).Elem()
		fill(t, rng, key, depth)
		value := reflect.New(v.Type().Elem()).Elem()
		fill(t, rng, value, depth)
		// a known and deliberate difference: gogo omits an empty bytes map
		// value, picobuf and the reference implementation encode it as a zero
		// length field. The only map that reaches signed bytes is
		// AddressedOrderLimit.tags, and the satellite no longer stores an
		// empty value there, see satellite/metainfo/endpoint_segment.go.
		if value.Kind() == reflect.Slice && value.Len() == 0 {
			value.SetBytes([]byte{1})
		}
		v.SetMapIndex(key, value)
	case reflect.Pointer:
		if depth >= maxDepth || rng.Intn(4) == 0 {
			return
		}
		v.Set(reflect.New(v.Type().Elem()))
		// now and then leave the pointed-to value at its defaults: the field
		// is still present, e.g. as an empty message or a zero timestamp.
		if rng.Intn(8) != 0 {
			fill(t, rng, v.Elem(), depth+1)
		}
		// a pointer to a nil slice, i.e. a custom typed field such as
		// ObjectListItem.stream_id declared without always_present, is the one
		// case the storj custom encoders cannot reproduce: gogo emits a zero
		// length field there. See storj/picobuf.go, it is unreachable in
		// signed messages and decodes to the same value either way.
		if v.Elem().Kind() == reflect.Slice && v.Elem().Len() == 0 {
			v.Elem().SetBytes([]byte{1})
		}
	case reflect.Struct:
		if depth > maxDepth {
			return
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			// XXX_unrecognized and friends are not part of the schema, and the
			// two implementations do not have the same set of them.
			if strings.HasPrefix(field.Name, "XXX_") {
				continue
			}
			// now and then leave an always_present message at its defaults,
			// which is still encoded, as an empty message.
			if field.Type.Kind() == reflect.Struct && field.Type != typeTime && rng.Intn(8) == 0 {
				continue
			}
			fill(t, rng, v.Field(i), depth)
		}
	case reflect.Interface:
		// a oneof field: pick one of the wrapper types. Both sides draw the
		// same index from the same seed and the generated lists are in the
		// same order, so both pick the same alternative.
		side := oneofPico
		if strings.HasSuffix(v.Type().PkgPath(), "/internal/gogopb") {
			side = oneofGogo
		}
		options := side[v.Type().Name()]
		if len(options) == 0 || depth >= maxDepth {
			return
		}
		wrapper := reflect.ValueOf(options[rng.Intn(len(options))]())
		inner := wrapper.Elem().Field(0)
		if inner.Kind() == reflect.Pointer {
			// a oneof alternative is always present when the field is set.
			inner.Set(reflect.New(inner.Type().Elem()))
			fill(t, rng, inner.Elem(), depth+1)
		} else {
			fill(t, rng, inner, depth+1)
		}
		v.Set(wrapper)
	default:
		t.Fatalf("unhandled kind %v for %v", v.Kind(), v.Type())
	}
}

func randBytes(rng *rand.Rand, n int) []byte {
	data := make([]byte, n)
	for i := range data {
		data[i] = byte(rng.Intn(256))
	}
	return data
}

func randString(rng *rand.Rand) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sb strings.Builder
	for i := 0; i < 1+rng.Intn(12); i++ {
		sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	return sb.String()
}
