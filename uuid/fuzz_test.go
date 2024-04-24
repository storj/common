// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

//go:build go1.18

package uuid_test

import (
	"testing"

	"storj.io/common/uuid"
)

func FuzzFromString(f *testing.F) {
	f.Add("")
	f.Add("6b\xff\xff\xff\u007f10-9dad-11d1-80b4-0c04fd430c8")
	f.Add("6ba7b810-$dad-11d1-80b4-0c04fd430c8")
	f.Add("6ba7b810-9dad-1 d1-.0b4-0c04fd430cA")
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("6ba7b810-9dad-11d1-80b4-0c04fd430c8")
	f.Add("6ba7b810-9dad-11d1-8>b4-0c04fd430c8")
	f.Add("6ba7b890-9dad-11d1-80b4-0c04fd430cF")
	f.Add("6Da7b890-9dad-11d1-80b4-0c04fd430cF")
	f.Add("6Da7b8D0-9dad-11d1-80b4-0c04CB430cF")
	f.Add("6Da7b8D0-9dad-11d1-80b4-0c04Cd430cF")
	f.Add("6Da7b8D0-9dad-11d1-80b4-0c04fd430cF")
	f.Add("6Da7b8D0-9dad-11d1-80b6-0c0FABSE\x10cF")
	f.Add("6Da7b8D0-9dad-11d1-80b6-0c0FALSE0cF")
	f.Add("6Da7b8D0-9dad-11d1-8FAB-0c0FABSE\x10cF")
	f.Add("6Da7b8D0-9dad-91d1-80cF-0c0FABSE\x10cF")
	f.Add("6Fa7b890-9dad-11d\xd7-80b4-0c04fd43/cb")
	f.Add("ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("EDa7b8D0-9dad-11d1-8FAB-0c0FABSE\x10cF")

	f.Fuzz(func(t *testing.T, data string) {
		_, _ = uuid.FromString(data)
	})
}
