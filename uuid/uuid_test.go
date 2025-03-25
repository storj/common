// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package uuid_test

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/testrand"
	"storj.io/common/uuid"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		s string
		u uuid.UUID
	}{
		// from RFC
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}},
		{"7d444840-9dc0-11d1-b245-5ffdce74fad2", uuid.UUID{0x7d, 0x44, 0x48, 0x40, 0x9d, 0xc0, 0x11, 0xd1, 0xb2, 0x45, 0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2}},
		// boundary cases
		{"00000000-0000-0000-0000-000000000000", uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"ffffffff-ffff-ffff-ffff-ffffffffffff", uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		// random
		{"0af42a8d-456e-4a68-af92-240413ffa492", uuid.UUID{0x0a, 0xf4, 0x2a, 0x8d, 0x45, 0x6e, 0x4a, 0x68, 0xaf, 0x92, 0x24, 0x04, 0x13, 0xff, 0xa4, 0x92}},
		{"99092ebc-ce69-4289-b496-bb58b376952e", uuid.UUID{0x99, 0x09, 0x2e, 0xbc, 0xce, 0x69, 0x42, 0x89, 0xb4, 0x96, 0xbb, 0x58, 0xb3, 0x76, 0x95, 0x2e}},
		{"b34df06f-be2b-4d62-a4f8-d1ae1f4fcfaa", uuid.UUID{0xb3, 0x4d, 0xf0, 0x6f, 0xbe, 0x2b, 0x4d, 0x62, 0xa4, 0xf8, 0xd1, 0xae, 0x1f, 0x4f, 0xcf, 0xaa}},
		{"755e114f-1e27-49ae-8cfc-873d9c6c6b10", uuid.UUID{0x75, 0x5e, 0x11, 0x4f, 0x1e, 0x27, 0x49, 0xae, 0x8c, 0xfc, 0x87, 0x3d, 0x9c, 0x6c, 0x6b, 0x10}},
		{"a11d29bd-5d1c-4a92-95cf-9d830c671811", uuid.UUID{0xa1, 0x1d, 0x29, 0xbd, 0x5d, 0x1c, 0x4a, 0x92, 0x95, 0xcf, 0x9d, 0x83, 0x0c, 0x67, 0x18, 0x11}},
		{"4fb0fc00-f584-4fca-a3eb-8a0c8709ef08", uuid.UUID{0x4f, 0xb0, 0xfc, 0x00, 0xf5, 0x84, 0x4f, 0xca, 0xa3, 0xeb, 0x8a, 0x0c, 0x87, 0x09, 0xef, 0x08}},
		{"b3e401e7-0137-4265-b874-ae2a79281026", uuid.UUID{0xb3, 0xe4, 0x01, 0xe7, 0x01, 0x37, 0x42, 0x65, 0xb8, 0x74, 0xae, 0x2a, 0x79, 0x28, 0x10, 0x26}},
		{"ade5b323-56e2-42e5-a347-51eb9d0e1272", uuid.UUID{0xad, 0xe5, 0xb3, 0x23, 0x56, 0xe2, 0x42, 0xe5, 0xa3, 0x47, 0x51, 0xeb, 0x9d, 0x0e, 0x12, 0x72}},
		// mixed case
		{"6Ba7B810-9dad-11D1-80B4-00C04Fd430c8", uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}},
		{"7d444840-9Dc0-11d1-b245-5ffdce74fad2", uuid.UUID{0x7d, 0x44, 0x48, 0x40, 0x9d, 0xc0, 0x11, 0xd1, 0xb2, 0x45, 0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2}},
		{"00000000-0000-0000-0000-000000000000", uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"ffffffff-fFFF-FFFF-FFff-ffffffffffff", uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}
	for _, test := range tests {
		x, err := uuid.FromString(test.s)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, test.u, x)
		assert.Equal(t, strings.ToLower(test.s), test.u.String())
	}

	tests = []struct {
		s string
		u uuid.UUID
	}{
		// from RFC
		{"6ba7b8109dad11d180b400c04fd430c8", uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}},
		{"7d4448409dc011d1b2455ffdce74fad2", uuid.UUID{0x7d, 0x44, 0x48, 0x40, 0x9d, 0xc0, 0x11, 0xd1, 0xb2, 0x45, 0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2}},
		// boundary cases
		{"00000000000000000000000000000000", uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"ffffffffffffffffffffffffffffffff", uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		// random
		{"0af42a8d456e4a68af92240413ffa492", uuid.UUID{0x0a, 0xf4, 0x2a, 0x8d, 0x45, 0x6e, 0x4a, 0x68, 0xaf, 0x92, 0x24, 0x04, 0x13, 0xff, 0xa4, 0x92}},
		{"99092ebcce694289b496bb58b376952e", uuid.UUID{0x99, 0x09, 0x2e, 0xbc, 0xce, 0x69, 0x42, 0x89, 0xb4, 0x96, 0xbb, 0x58, 0xb3, 0x76, 0x95, 0x2e}},
		{"b34df06fbe2b4d62a4f8d1ae1f4fcfaa", uuid.UUID{0xb3, 0x4d, 0xf0, 0x6f, 0xbe, 0x2b, 0x4d, 0x62, 0xa4, 0xf8, 0xd1, 0xae, 0x1f, 0x4f, 0xcf, 0xaa}},
		{"755e114f1e2749ae8cfc873d9c6c6b10", uuid.UUID{0x75, 0x5e, 0x11, 0x4f, 0x1e, 0x27, 0x49, 0xae, 0x8c, 0xfc, 0x87, 0x3d, 0x9c, 0x6c, 0x6b, 0x10}},
		{"a11d29bd5d1c4a9295cf9d830c671811", uuid.UUID{0xa1, 0x1d, 0x29, 0xbd, 0x5d, 0x1c, 0x4a, 0x92, 0x95, 0xcf, 0x9d, 0x83, 0x0c, 0x67, 0x18, 0x11}},
		{"4fb0fc00f5844fcaa3eb8a0c8709ef08", uuid.UUID{0x4f, 0xb0, 0xfc, 0x00, 0xf5, 0x84, 0x4f, 0xca, 0xa3, 0xeb, 0x8a, 0x0c, 0x87, 0x09, 0xef, 0x08}},
		{"b3e401e701374265b874ae2a79281026", uuid.UUID{0xb3, 0xe4, 0x01, 0xe7, 0x01, 0x37, 0x42, 0x65, 0xb8, 0x74, 0xae, 0x2a, 0x79, 0x28, 0x10, 0x26}},
		{"ade5b32356e242e5a34751eb9d0e1272", uuid.UUID{0xad, 0xe5, 0xb3, 0x23, 0x56, 0xe2, 0x42, 0xe5, 0xa3, 0x47, 0x51, 0xeb, 0x9d, 0x0e, 0x12, 0x72}},
		// mixed case
		{"6Ba7B8109dad11D180B400C04Fd430c8", uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}},
		{"7d4448409Dc011d1b2455ffdce74fad2", uuid.UUID{0x7d, 0x44, 0x48, 0x40, 0x9d, 0xc0, 0x11, 0xd1, 0xb2, 0x45, 0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2}},
		{"00000000000000000000000000000000", uuid.UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"fffffffffFFFFFFFFFffffffffffffff", uuid.UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	}
	for _, test := range tests {
		x, err := uuid.FromString(test.s)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, test.u, x)

		var expected string
		{
			s := strings.ToLower(test.s)
			expected = s[0:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:32]
		}

		assert.Equal(t, expected, test.u.String())
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x, err := uuid.New()
		require.NoError(t, err)
		require.False(t, x.IsZero())

		parsed, err := uuid.FromString(x.String())
		require.NoError(t, err)
		require.Equal(t, x, parsed)
	}
}

func TestJSON(t *testing.T) {
	type example struct {
		A uuid.UUID  `json:"A"`
		B *uuid.UUID `json:"B"`
	}
	var x example
	x.A = uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	x.B = &uuid.UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	expected := `{"A":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","B":"6ba7b810-9dad-11d1-80b4-00c04fd430c8"}`

	data, err := json.Marshal(x)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))

	var b example
	err = json.Unmarshal([]byte(expected), &b)
	require.NoError(t, err)
	require.Equal(t, x, b)

	expected = `{"A":"6ba7b8109dad11d180b400c04fd430c8","B":"6ba7b8109dad11d180b400c04fd430c8"}`
	err = json.Unmarshal([]byte(expected), &b)
	require.NoError(t, err)
	require.Equal(t, x, b)
}

func TestLess(t *testing.T) {
	for k := 0; k < len(uuid.UUID{}); k++ {
		var a, b uuid.UUID
		a[k], b[k] = 1, 2
		require.True(t, a.Less(b))
	}

	for k := 0; k < 100; k++ {
		var x, y uuid.UUID
		a, b := testrand.Int63n(math.MaxInt64), testrand.Int63n(math.MaxInt64)
		binary.BigEndian.PutUint64(x[:], uint64(a))
		binary.BigEndian.PutUint64(y[:], uint64(b))
		require.Equal(t, a < b, x.Less(y))
	}
}

func TestSort(t *testing.T) {
	uuid1, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	uuid2, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	uuid3, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	emptyAsc := []uuid.UUID{}
	uuid.SortAscending(emptyAsc)
	require.Equal(t, []uuid.UUID{}, emptyAsc)

	oneAsc := []uuid.UUID{uuid1}
	uuid.SortAscending(oneAsc)
	require.Equal(t, []uuid.UUID{uuid1}, oneAsc)

	threeAscToAsc := []uuid.UUID{uuid1, uuid2, uuid3}
	uuid.SortAscending(threeAscToAsc)
	require.Equal(t, []uuid.UUID{uuid1, uuid2, uuid3}, threeAscToAsc)

	threeDescToAsc := []uuid.UUID{uuid3, uuid2, uuid1}
	uuid.SortAscending(threeDescToAsc)
	require.Equal(t, []uuid.UUID{uuid1, uuid2, uuid3}, threeDescToAsc)
}

func TestMarshal(t *testing.T) {
	expectedUUID := testrand.UUID()
	uuidBytes, err := expectedUUID.Marshal()
	require.NoError(t, err)

	uuid := uuid.UUID{}
	err = uuid.Unmarshal(uuidBytes)
	require.NoError(t, err)
	require.Equal(t, expectedUUID, uuid)
	require.Equal(t, expectedUUID.Bytes(), uuid.Bytes())
	require.Equal(t, expectedUUID.Size(), uuid.Size())

	uuidBytesTo := make([]byte, len(expectedUUID))
	n, err := expectedUUID.MarshalTo(uuidBytesTo)
	require.NoError(t, err)
	require.Equal(t, len(expectedUUID), n)
	require.Equal(t, uuidBytes, uuidBytesTo)
}

func TestCompare(t *testing.T) {
	var a uuid.UUID
	require.Equal(t, 0, a.Compare(a)) //nolint:gocritic

	for k := 0; k < len(uuid.UUID{}); k++ {
		var a, b uuid.UUID
		a[k], b[k] = 1, 2
		require.Equal(t, 0, a.Compare(a)) //nolint:gocritic
		require.Equal(t, 0, b.Compare(b)) //nolint:gocritic
		require.Equal(t, -1, a.Compare(b))
		require.Equal(t, 1, b.Compare(a))
	}

	for k := 0; k < 100; k++ {
		var x, y uuid.UUID
		a, b := testrand.Int63n(math.MaxInt64), testrand.Int63n(math.MaxInt64)
		binary.BigEndian.PutUint64(x[:], uint64(a))
		binary.BigEndian.PutUint64(y[:], uint64(b))
		require.Equal(t, comp(a, b), x.Compare(y))
	}
}

func comp(a, b int64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func BenchmarkLess(b *testing.B) {
	a := testrand.UUID()
	b.Run("Same", func(b *testing.B) {
		total := 0
		x, y := a, a
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("First", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[0]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("Middle", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)/2]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})

	b.Run("Last", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)-1]++
		for k := 0; k < b.N; k++ {
			total += btoi(x.Less(y))
		}
		runtime.KeepAlive(total)
	})
}

func btoi(v bool) int {
	if v {
		return 1
	}
	return 0
}

func BenchmarkCompare(b *testing.B) {
	a := testrand.UUID()
	b.Run("Same", func(b *testing.B) {
		total := 0
		x, y := a, a
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("First", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[0]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("Middle", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)/2]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})

	b.Run("Last", func(b *testing.B) {
		total := 0
		x, y := a, a
		y[len(y)-1]++
		for k := 0; k < b.N; k++ {
			total += x.Compare(y)
		}
		runtime.KeepAlive(total)
	})
}
