// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information.

package process_test

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"storj.io/common/process"
	"storj.io/common/testcontext"
)

func TestSaveConfig_CustomFlagName(t *testing.T) {
	cmd := &cobra.Command{RunE: func(cmd *cobra.Command, args []string) error { return nil }}

	var config struct {
		Flag string `flagname:"flag-name" default:"value"`
	}
	process.Bind(cmd, &config)

	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	configFile := ctx.File(process.DefaultCfgFilename)

	require.NoError(t, process.SaveConfig(cmd, configFile))

	configContents, err := os.ReadFile(configFile)
	require.NoError(t, err)
	require.Equal(t, "# flag-name: value\n", string(configContents))
}
