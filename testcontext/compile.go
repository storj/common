// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package testcontext

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// Compile compiles the specified package and returns the executable name.
func (ctx *Context) Compile(pkg string, preArgs ...string) string {
	ctx.test.Helper()

	var binName string
	if pkg == "" {
		dir, _ := os.Getwd()
		binName = path.Base(dir)
	} else {
		binName = path.Base(pkg)
	}

	exe := ctx.File("build", binName+".exe")

	args := append([]string{"build"}, preArgs...)
	if raceEnabled {
		args = append(args, "-race")
	}
	args = append(args, "-tags=unittest")
	args = append(args, "-o", exe, pkg)

	/* #nosec G204 */ // This package is only used for test
	cmd := exec.Command("go", args...)
	ctx.test.Log("exec:", cmd.Args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		ctx.test.Error(string(out))
		ctx.test.Fatal(err)
	}

	return exe
}

// CompileAt compiles the specified package and returns the executable name.
func (ctx *Context) CompileAt(workDir, pkg string, preArgs ...string) string {
	ctx.test.Helper()

	var binName string
	if pkg == "" {
		dir, _ := os.Getwd()
		binName = path.Base(dir)
	} else {
		binName = path.Base(pkg)
	}

	if absDir, err := filepath.Abs(workDir); err == nil {
		workDir = absDir
	} else {
		ctx.test.Fatal(err)
	}

	exe := ctx.File("build", binName+".exe")

	args := append([]string{"build"}, preArgs...)
	if raceEnabled {
		args = append(args, "-race")
	}
	args = append(args, "-tags=unittest")
	args = append(args, "-o", exe, pkg)

	/* #nosec G204 */ // This package is only used for test
	cmd := exec.Command("go", args...)
	cmd.Dir = workDir
	ctx.test.Log("exec:", cmd.Args, "dir:", workDir)

	out, err := cmd.CombinedOutput()
	if err != nil {
		ctx.test.Error(string(out))
		ctx.test.Fatal(err)
	}

	return exe
}

// CompileWithLDFlagsX compiles the specified package with the -ldflags flag set to
// "-s -w [-X <key>=<value>,...]" given the passed map and returns the executable name.
func (ctx *Context) CompileWithLDFlagsX(pkg string, ldFlagsX map[string]string) string {
	ctx.test.Helper()

	var ldFlags = "-s -w"
	for key, value := range ldFlagsX {
		ldFlags += (" -X " + key + "=" + value)
	}

	return ctx.Compile(pkg, "-ldflags", ldFlags)
}
