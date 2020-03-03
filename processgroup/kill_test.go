// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information

package processgroup_test

import (
	"io/ioutil"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"storj.io/common/processgroup"
	"storj.io/common/testcontext"
)

func TestProcessGroup(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	source := ctx.File("main.go")
	binary := ctx.File("main.exe")
	err := ioutil.WriteFile(source, []byte(code), 0644)
	require.NoError(t, err)

	{
		cmd := exec.Command("go", "build", "-o", binary, source)
		cmd.Dir = ctx.Dir()

		_, err := cmd.CombinedOutput()
		require.NoError(t, err)
	}

	{
		cmd := exec.Command(binary)
		cmd.Dir = ctx.Dir()
		cmd.Stdout, cmd.Stderr = ioutil.Discard, ioutil.Discard
		processgroup.Setup(cmd)

		started := time.Now()
		err := cmd.Start()
		require.NoError(t, err)
		processgroup.Kill(cmd)

		_ = cmd.Wait() // since we kill it, we might get an error
		duration := time.Since(started)

		require.Truef(t, duration < 10*time.Second, "completed in %s", duration)
	}
}

const code = `package main

import "time"

func main() {
	time.Sleep(20*time.Second)
}
`
