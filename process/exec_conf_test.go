// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
)

func setenv(key, value string) func() {
	old := os.Getenv(key)
	_ = os.Setenv(key, value)
	return func() { _ = os.Setenv(key, old) }
}

func setargs(value []string) func() {
	old := os.Args
	os.Args = value
	return func() { os.Args = old }
}

var testZ = flag.Int("z", 0, "z flag (stdlib)")

func TestExec_PropagatesSettings(t *testing.T) {
	// Set up a command that does nothing.
	cmd := &cobra.Command{RunE: func(cmd *cobra.Command, args []string) error { return nil }}

	// Define a config struct and some flags.
	var config struct {
		X int `default:"0"`
	}
	Bind(cmd, &config)
	y := cmd.Flags().Int("y", 0, "y flag (command)")

	// Set some environment variables for viper.
	defer setenv("STORJ_X", "1")()
	defer setenv("STORJ_Y", "2")()
	defer setenv("STORJ_Z", "3")()

	// Run the command through the exec call.
	Exec(cmd)

	// Check that the variables are now bound.
	require.Equal(t, 1, config.X)
	require.Equal(t, 2, *y)
	require.Equal(t, 3, *testZ)
}

func TestExec_InvalidValues(t *testing.T) {

	// previous test adds golang specific flags (see Exec) which pollutes the tests here.
	oldPflagCommandline := pflag.CommandLine
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	defer func() {
		pflag.CommandLine = oldPflagCommandline
	}()

	noEnv := map[string]string{}
	var noFlags []string
	failOnError := func(failOnError bool) func(options *ExecOptions) {
		return func(options *ExecOptions) {
			options.FailOnValueError = failOnError
		}
	}
	runCommand := func(t *testing.T, flags []string, env map[string]string, configFileContent string, execOptions func(options *ExecOptions)) error {
		// Set up a command that does nothing.
		cmd := &cobra.Command{RunE: func(cmd *cobra.Command, args []string) error { return nil }}

		// Define a config struct and some flags.
		var config struct {
			X storj.NodeURLs `default:"0"`
		}
		args := []string{"test"}
		args = append(args, flags...)
		defer setargs(args)()

		for k, v := range env {
			defer setenv(k, v)()
		}
		Bind(cmd, &config)

		configFile := filepath.Join(t.TempDir(), t.Name()+"config.yaml")
		_ = os.MkdirAll(filepath.Dir(configFile), 0755)
		err := os.WriteFile(configFile, []byte(configFileContent), 0644)
		require.NoError(t, err)

		execOpts := &ExecOptions{
			LoadConfig: func(cmd *cobra.Command, vip *viper.Viper) error {
				vip.SetConfigFile(configFile)
				return vip.ReadInConfig()
			},
		}

		execOptions(execOpts)
		cleanup(cmd, execOpts)

		return cmd.Execute()

	}
	for _, failOnValueProblem := range []bool{false, true} {
		t.Run(fmt.Sprintf("strict_value_%v", failOnValueProblem), func(t *testing.T) {
			t.Run("Set by args", func(t *testing.T) {
				t.Run("mistyped value", func(t *testing.T) {
					err := runCommand(t, []string{"--x", "not-a-nodeid@localhost"}, noEnv, "", failOnError(failOnValueProblem))

					// flags are always failing, as before, because we can be sure that's an error.
					require.Error(t, err)

				})
				t.Run("mistyped key", func(t *testing.T) {
					err := runCommand(t, []string{"--xy", "asd"}, noEnv, "", failOnError(failOnValueProblem))

					// flags are always failing, as before.
					require.Error(t, err)
				})

			})

			t.Run("Set by env", func(t *testing.T) {
				t.Run("mistyped value", func(t *testing.T) {
					err := runCommand(t, noFlags, map[string]string{"STORJ_X": "not-a-node-id@localhost"}, "", failOnError(failOnValueProblem))

					if failOnValueProblem {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}
				})
				t.Run("mistyped key", func(t *testing.T) {
					err := runCommand(t, noFlags, map[string]string{"STORJ_XY": "not-a-node-id@localhost"}, "", failOnError(failOnValueProblem))

					// It might be a config value for other service type. We couldn't fail without checking others.
					require.NoError(t, err)

				})

			})

			t.Run("Set by file", func(t *testing.T) {
				t.Run("mistyped value", func(t *testing.T) {
					err := runCommand(t, noFlags, noEnv, "x: not-a-node-id@localhost", failOnError(failOnValueProblem))

					if failOnValueProblem {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}
				})
				t.Run("mistyped key", func(t *testing.T) {
					err := runCommand(t, noFlags, noEnv, "xy: localhost", failOnError(failOnValueProblem))

					// Might be a config value for an other service.
					require.NoError(t, err)
				})
			})
		})

	}

}

func TestHidden(t *testing.T) {
	// Set up a command that does nothing.
	cmd := &cobra.Command{RunE: func(cmd *cobra.Command, args []string) error { return nil }}

	// Define a config struct with a hidden field.
	var config struct {
		W int `default:"0" hidden:"false"`
		X int `default:"0" hidden:"true"`
		Y int `releaseDefault:"1" devDefault:"0" hidden:"true"`
		Z int `default:"1"`
	}
	Bind(cmd, &config)

	// Setup test config file
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	testConfigFile := ctx.File("testconfig.yaml")

	// Run the command through the exec call.
	Exec(cmd)

	// Ensure that the file saves only the necessary data.
	err := SaveConfig(cmd, testConfigFile)
	require.NoError(t, err)

	/* #nosec G304 */ // The file path is generated by testcontext so there isn't
	// any security flaw of a file inclusion via a variable
	actualConfigFile, err := os.ReadFile(testConfigFile)
	require.NoError(t, err)

	expectedConfigW := "# w: 0"
	expectedConfigZ := "# z: 1"
	require.Contains(t, string(actualConfigFile), expectedConfigW)
	require.Contains(t, string(actualConfigFile), expectedConfigZ)
	require.NotContains(t, string(actualConfigFile), "# y: ")
	require.NotContains(t, string(actualConfigFile), "# x: ")
}
