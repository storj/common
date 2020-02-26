// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	grpcdir = flag.String("grpc-dir", "pbgrpc", "directory for grpc client and servers")
	mainpkg = flag.String("pkg", "storj.io/common/pb", "main package name")
	protoc  = flag.String("protoc", "protoc", "protoc compiler")
)

var ignoreProto = map[string]bool{
	"gogo.proto": true,
}

func ignore(files []string) []string {
	xs := []string{}
	for _, file := range files {
		if !ignoreProto[file] {
			xs = append(xs, file)
		}
	}
	return xs
}

func main() {
	flag.Parse()

	// TODO: protolock

	{
		out, err := exec.Command("go", "run", "../scripts/protobuf.go", "lint").CombinedOutput()
		fmt.Println(string(out))
		check(err)
	}

	// create grpcdir if it doesn't exist
	_ = os.Mkdir(*grpcdir, 0644)

	{
		// cleanup previous files
		grpcfiles, err := filepath.Glob(filepath.Join(*grpcdir, "*.pb.go"))
		check(err)
		localfiles, err := filepath.Glob("*.pb.go")
		check(err)

		all := []string{}
		all = append(all, grpcfiles...)
		all = append(all, localfiles...)
		for _, match := range all {
			_ = os.Remove(match)
		}
	}

	{
		protofiles, err := filepath.Glob("*.proto")
		check(err)

		protofiles = ignore(protofiles)

		args := []string{
			"--drpc_out=plugins=grpc+drpc:.",
			"-I=.",
		}
		args = append(args, protofiles...)

		// generate new code
		cmd := exec.Command(*protoc, args...)
		fmt.Println(strings.Join(cmd.Args, " "))
		out, err := cmd.CombinedOutput()
		fmt.Println(string(out))
		check(err)
	}

	{
		// split generated files
		files, err := filepath.Glob("*.pb.go")
		check(err)
		for _, file := range files {
			split(file)
		}
	}

	{
		// format code to get rid of extra imports
		out, err := exec.Command("goimports", "-local", "storj.io", "-w", ".").CombinedOutput()
		fmt.Println(string(out))
		check(err)
	}
}

const drpcEndTag = "// --- DRPC END ---\n"

// split moves grpc part of the code to grpcdir folder.
func split(file string) {
	data, err := ioutil.ReadFile(file)
	check(err)

	source := string(data)

	// first ')' is the closing of import block
	importEnd := strings.IndexByte(source, ')')

	// find drpc tag
	drpcEnd := strings.Index(source, drpcEndTag)
	if drpcEnd < 0 {
		// no service code
		return
	}
	drpcEnd += len(drpcEndTag)

	// create source for grpc
	grpc := ""
	grpc += source[:importEnd] + fmt.Sprintf("\t. %q\n)\n", *mainpkg)
	grpc += source[drpcEnd:]

	grpc = strings.Replace(grpc, "package pb", "package "+*grpcdir, -1)

	err = ioutil.WriteFile(filepath.Join(*grpcdir, file), []byte(grpc), 0644)
	check(err)

	// rewrite original file without grpc
	err = ioutil.WriteFile(file, []byte(source[:drpcEnd]), 0644)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
