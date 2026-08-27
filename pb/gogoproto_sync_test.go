// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

package pb_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGogoProtoInSync verifies that the .proto copies in internal/gogopb
// declare the same schema as the ones here, ignoring comments, imports, file
// options and field options, which is where the pico and the gogoproto
// annotations differ. Without it a change made only on one side would leave
// TestGogoCompatibility comparing against a stale definition.
func TestGogoProtoInSync(t *testing.T) {
	files, err := filepath.Glob("*.proto")
	require.NoError(t, err)

	gogoFiles, err := filepath.Glob(filepath.Join("internal", "gogopb", "*.proto"))
	require.NoError(t, err)

	var names, gogoNames []string
	for _, file := range files {
		if file != "pico.proto" {
			names = append(names, file)
		}
	}
	for _, file := range gogoFiles {
		if name := filepath.Base(file); name != "gogo.proto" {
			gogoNames = append(gogoNames, name)
		}
	}
	require.Equal(t, names, gogoNames, "internal/gogopb must have a copy of every .proto")

	for _, name := range names {
		pico, err := os.ReadFile(name)
		require.NoError(t, err)
		gogo, err := os.ReadFile(filepath.Join("internal", "gogopb", name))
		require.NoError(t, err)
		require.Equal(t, protoSchema(string(pico)), protoSchema(string(gogo)), "%s differs from internal/gogopb/%s", name, name)
	}
}

var (
	rxProtoComment = regexp.MustCompile(`(?s)//[^\n]*|/\*.*?\*/`)
	rxProtoHeader  = regexp.MustCompile(`(?m)^\s*(import|option)\s[^;]*;`)
	rxFieldOptions = regexp.MustCompile(`\[[^\]]*\]`)
)

// protoSchema strips everything but the declarations from a .proto file and
// returns them one statement per line.
func protoSchema(source string) string {
	source = rxProtoComment.ReplaceAllString(source, "")
	source = rxProtoHeader.ReplaceAllString(source, "")
	source = rxFieldOptions.ReplaceAllString(source, "")
	source = strings.Join(strings.Fields(source), " ")
	for _, sep := range []string{";", "{", "}"} {
		source = strings.ReplaceAll(source, " "+sep, sep)
		source = strings.ReplaceAll(source, sep, sep+"\n")
	}
	return source
}
