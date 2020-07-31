// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package memory_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/memory"
)

const (
	eib = 1 << 60
	pib = 1 << 50
	tib = 1 << 40
	gib = 1 << 30
	mib = 1 << 20
	kib = 1 << 10
	eb  = 1e18
	pb  = 1e15
	tb  = 1e12
	gb  = 1e9
	mb  = 1e6
	kb  = 1e3
)

func TestBase2Size(t *testing.T) {
	var tests = []struct {
		size memory.Size
		text string
	}{
		// basics
		{1 * eib, "1.0 EiB"},
		{1 * pib, "1.0 PiB"},
		{1 * tib, "1.0 TiB"},
		{1 * gib, "1.0 GiB"},
		{1 * mib, "1.0 MiB"},
		{1 * kib, "1.0 KiB"},
		{1, "1 B"},
		// complicated
		{68 * tib, "68.0 TiB"},
		{256 * mib, "256.0 MiB"},
		{500, "500 B"},
		{5, "5 B"},
		{1, "1 B"},
		{0, "0 B"},
	}

	for i, test := range tests {
		if test.size.String() != test.text {
			t.Errorf("%d. invalid text got %v expected %v", i, test.size.String(), test.text)
		}

		var size memory.Size
		err := size.Set(test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}

		// test negative
		err = size.Set("-" + test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if -test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}
	}
}

func TestBase10Size(t *testing.T) {
	var tests = []struct {
		size memory.Size
		text string
	}{
		// basics
		{1 * pb, "1.00 PB"},
		{1 * eb, "1.00 EB"},
		{1 * tb, "1.00 TB"},
		{1 * gb, "1.00 GB"},
		{1 * mb, "1.00 MB"},
		{1 * kb, "1.00 KB"},
		{1, "1 B"},
		// complicated
		{68 * tb, "68.00 TB"},
		{256 * mb, "256.00 MB"},
		{500, "500 B"},
		{5, "5 B"},
		{1, "1 B"},
		{0, "0 B"},
	}

	for i, test := range tests {
		if test.size.Base10String() != test.text {
			t.Errorf("%d. invalid text got %v expected %v", i, test.size.String(), test.text)
		}

		var size memory.Size
		err := size.Set(test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}

		// test negative
		err = size.Set("-" + test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if -test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}
	}
}

func TestParse(t *testing.T) {
	var tests = []struct {
		size memory.Size
		text string
	}{
		// case insensitivity
		{1 * tb, "1.00TB"},
		{1 * gb, "1.00gB"},
		{1 * mb, "1.00Mb"},
		{1 * kb, "1.00kb"},
		{1, "1.00"},
		{1 * tb, "1.0 TB"},
		{1 * gb, "1.0 gB"},
		{1 * mb, "1.0 Mb"},
		{1 * kb, "1.0 kb"},
		{1 * kib, "1.0kib"},
		{1 * pib, "1.0pib"},
		{1 * eib, "1.0eib"},
		{1, "1.00"},
		// without B suffix
		{1 * tb, "1.00T"},
		{1 * gb, "1.00g"},
		{1 * mb, "1.00M"},
		{1 * kb, "1.00k"},
	}

	for i, test := range tests {
		var size memory.Size
		err := size.Set(test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}

		// test negative
		err = size.Set("-" + test.text)
		if err != nil {
			t.Errorf("%d. got error %v", i, err)
		}
		if -test.size != size {
			t.Errorf("%d. invalid size got %d expected %d", i, size, test.size)
		}
	}
}

func TestInvalidParse(t *testing.T) {
	var tests = []string{
		"1.0Q",
		"1.0QB",
		"z1.0KB",
		"z1.0Q",
		"1.0zQ",
		"1.0zQB",
	}

	for i, test := range tests {
		var size memory.Size
		err := size.Set(test)
		if err == nil {
			t.Errorf("%d. didn't get error", i)
		}
	}
}

func TestJSON(t *testing.T) {
	var input struct {
		Value memory.Size `json:"value"`
	}
	input.Value = memory.TiB

	data, err := json.Marshal(input)
	require.NoError(t, err)
	require.Equal(t, `{"value":"1.0 TiB"}`, string(data))

	err = json.Unmarshal([]byte(`{"value" : "1GB"}`), &input)
	require.NoError(t, err)
	require.Equal(t, memory.GB, input.Value)
}
