// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// +build gofuzz

package useragent

// To run fuzzing tests:
//
// clone github.com/storj/fuzz-corpus
//
// Install fuzzing tools:
//   GO111MODULE=off go get github.com/dvyukov/go-fuzz/...
//
// Build binaries:
//   go-fuzz-build .
//
// Run with test corpus:
//   go-fuzz -bin useragent-fuzz.zip -workdir $FUZZCORPUS/useragent/testdata

// Fuzz implements a simple fuzz test for uuid.Parse.
func Fuzz(data []byte) int {
	_, err := ParseEntries(data)
	if err != nil {
		return 0
	}
	return 1
}
