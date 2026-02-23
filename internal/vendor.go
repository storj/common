// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// This vendor vendors crypto/sha512 block implementation from std Go.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:generate go run .

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	goroot, err := goEnvGOROOT()
	if err != nil {
		return err
	}

	// As of Go 1.24, the sha512 block implementation moved from
	// crypto/sha512 to crypto/internal/fips140/sha512.
	sha512block := filepath.Join(goroot, "src", "crypto", "internal", "fips140", "sha512", "sha512block*")

	matches, err := filepath.Glob(sha512block)
	if err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("no files matching %s", sha512block)
	}

	// Skip the _asm directory.
	var filtered []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}
		filtered = append(filtered, match)
	}
	matches = filtered

	rxImplRegister := regexp.MustCompile(`(?ms)^func init\(\) \{[^}]*impl\.Register[^}]*\}\n`)
	rxImplImport := regexp.MustCompile(`(?m)"crypto/internal/impl"\n?\t?`)
	rxEmptyImport := regexp.MustCompile(`(?m)^import \(\s*\)\n`)

	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			return err
		}

		data = bytes.ReplaceAll(data, []byte(`package sha512`), []byte(`package hmacsha512`))
		data = bytes.ReplaceAll(data,
			[]byte(`// license that can be found in the LICENSE file.`),
			[]byte(`// license that can be found in the GO_LICENSE file.`),
		)

		// Replace internal FIPS cpu package with golang.org/x/sys/cpu.
		data = bytes.ReplaceAll(data,
			[]byte(`"crypto/internal/fips140deps/cpu"`),
			[]byte(`"golang.org/x/sys/cpu"`))

		// Translate cpu feature variable names from fips140deps style to x/sys/cpu style.
		data = bytes.ReplaceAll(data, []byte(`cpu.X86HasAVX2`), []byte(`cpu.X86.HasAVX2`))
		data = bytes.ReplaceAll(data, []byte(`cpu.X86HasAVX`), []byte(`cpu.X86.HasAVX`))
		data = bytes.ReplaceAll(data, []byte(`cpu.X86HasBMI2`), []byte(`cpu.X86.HasBMI2`))
		data = bytes.ReplaceAll(data, []byte(`cpu.ARM64HasSHA512`), []byte(`cpu.ARM64.HasSHA512`))
		data = bytes.ReplaceAll(data, []byte(`cpu.S390XHasSHA512`), []byte(`cpu.S390X.HasSHA512`))

		// The vendored package uses unexported "digest" instead of "Digest".
		data = bytes.ReplaceAll(data, []byte(`*Digest`), []byte(`*digest`))

		// Remove impl.Register calls and the crypto/internal/impl import.
		data = rxImplRegister.ReplaceAll(data, nil)
		data = rxImplImport.ReplaceAll(data, nil)

		// Remove purego from build constraints — not relevant for external packages.
		data = bytes.ReplaceAll(data, []byte("//go:build !purego\n"), []byte(""))
		data = bytes.ReplaceAll(data, []byte(" && !purego"), []byte(""))
		data = bytes.ReplaceAll(data, []byte(" || purego"), []byte(""))

		// Handle ppc64x: replace godebug-based feature detection with always-on.
		data = bytes.ReplaceAll(data,
			[]byte(`"crypto/internal/fips140deps/godebug"`),
			[]byte{})
		data = bytes.ReplaceAll(data,
			[]byte(`var ppc64sha512 = godebug.Value("#ppc64sha512") != "off"`),
			[]byte(`var ppc64sha512 = true`))

		// Clean up empty import blocks left after removing imports.
		data = rxEmptyImport.ReplaceAll(data, nil)

		err = os.WriteFile(filepath.Join("hmacsha512", filepath.Base(match)), data, 0o644)
		if err != nil {
			return err
		}
	}

	version, err := os.ReadFile(filepath.Join(goroot, "VERSION"))
	if err != nil {
		return err
	}

	// VERSION file may contain multiple lines (e.g. "go1.26.0\ntime ..."),
	// only use the first line.
	versionLine, _, _ := strings.Cut(strings.TrimSpace(string(version)), "\n")
	doc := strings.ReplaceAll(doc, "{{.Version}}", versionLine)
	err = os.WriteFile(filepath.Join("hmacsha512", "doc.go"), []byte(doc), 0o644)
	if err != nil {
		return err
	}

	out, err := exec.Command("go", "fmt", "./hmacsha512").CombinedOutput()
	if err != nil {
		return fmt.Errorf("go fmt ./hmacsha512: %w\n%s", err, out)
	}

	return nil
}

func goEnvGOROOT() (string, error) {
	out, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go env GOROOT: %w\n%s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

const doc = `// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// Package hmacsha512 contains an inlined an optimized version of hmac+sha512.
// Unfortunately, this requires exposing some of the details from crypto/sha512.
package hmacsha512

// Currently vendored crypto/sha512 version is {{.Version}}
`
