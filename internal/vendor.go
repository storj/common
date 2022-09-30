// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// This vendor vendors crypto/sha512 block implementation from std Go.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
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
	sha512block := filepath.Join(runtime.GOROOT(), "src", "crypto", "sha512", "sha512block*")

	matches, err := filepath.Glob(sha512block)
	if err != nil {
		return err
	}

	rxGoBuild := regexp.MustCompile("(?m)^//go:build .*$")

	for _, match := range matches {
		data, err := ioutil.ReadFile(match)
		if err != nil {
			return err
		}

		data = bytes.ReplaceAll(data, []byte(`package sha512`), []byte(`package hmacsha512`))
		data = bytes.ReplaceAll(data, []byte(`import "internal/cpu"`), []byte(`import "golang.org/x/sys/cpu"`))
		data = bytes.ReplaceAll(data,
			[]byte(`// license that can be found in the LICENSE file.`),
			[]byte(`// license that can be found in the GO_LICENSE file.`),
		)

		// Ensure we preseve the old build tags to be compatible with old Go version.
		data = rxGoBuild.ReplaceAll(data, []byte("$0\n// +build stub\n"))

		err = ioutil.WriteFile(filepath.Join("hmacsha512", filepath.Base(match)), data, 0755)
		if err != nil {
			return err
		}
	}

	version, err := ioutil.ReadFile(filepath.Join(runtime.GOROOT(), "VERSION"))
	if err != nil {
		return err
	}

	doc := strings.ReplaceAll(doc, "{{.Version}}", strings.TrimSpace(string(version)))
	err = ioutil.WriteFile(filepath.Join("hmacsha512", "doc.go"), []byte(doc), 0755)
	if err != nil {
		return err
	}

	_, err = exec.Command("go", "fmt", "./hmacsha512").CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

const doc = `// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// Package hmacsha512 contains an inlined an optimized version of hmac+sha512.
// Unfortunately, this requires exposing some of the details from crypto/sha512.
package hmacsha512

// Currently vendored crypto/sha512 version is {{.Version}}
`
