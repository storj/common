// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spacemonkeygo/monkit/v3/collect"
	"github.com/spacemonkeygo/monkit/v3/present"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zeebo/errs"
	"github.com/zeebo/structs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"storj.io/common/cfgstruct"
	"storj.io/common/version"
)

// DefaultCfgFilename is the default filename used for storing a configuration.
const DefaultCfgFilename = "config.yaml"

var (
	mon = monkit.Package()

	commandMtx   sync.Mutex
	contexts     = map[*cobra.Command]context.Context{}
	cancels      = map[*cobra.Command]context.CancelFunc{}
	configs      = map[*cobra.Command][]interface{}{}
	vipers       = map[*cobra.Command]*viper.Viper{}
	atomicLevels = map[*cobra.Command]*zap.AtomicLevel{}
)

// Bind sets flags on a command that match the configuration struct
// 'config'. It ensures that the config has all of the values loaded into it
// when the command runs.
func Bind(cmd *cobra.Command, config interface{}, opts ...cfgstruct.BindOpt) {
	commandMtx.Lock()
	defer commandMtx.Unlock()

	cfgstruct.Bind(cmd.Flags(), config, opts...)
	configs[cmd] = append(configs[cmd], config)
}

// Exec runs a Cobra command. If a "config-dir" flag is defined it will be parsed
// and loaded using viper.
func Exec(cmd *cobra.Command) {
	ExecWithCustomConfig(cmd, true, LoadConfig)
}

// ExecCustomDebug runs default configuration except the default debug is disabled.
func ExecCustomDebug(cmd *cobra.Command) {
	ExecWithCustomConfig(cmd, false, LoadConfig)
}

// ExecWithCustomConfig runs a Cobra command. Custom configuration can be loaded.
func ExecWithCustomConfig(cmd *cobra.Command, debugEnabled bool, loadConfig func(cmd *cobra.Command, vip *viper.Viper) error) {
	ExecWithCustomConfigAndLogger(cmd, debugEnabled, loadConfig, nil)
}

// ExecWithCustomConfigAndLogger runs a Cobra command with a custom logger. Custom configuration can be loaded.
func ExecWithCustomConfigAndLogger(cmd *cobra.Command, debugEnabled bool, loadConfig func(cmd *cobra.Command, vip *viper.Viper) error, loggerFactory func(*zap.Logger) *zap.Logger) {
	ExecWithCustomOptions(cmd, ExecOptions{
		InitDefaultDebugServer: debugEnabled,
		InitProfiler:           true,
		InitTracing:            true,
		LoadConfig:             loadConfig,
		LoggerFactory:          loggerFactory,
	})
}

// ExecOptions contains options for ExecWithCustomOptions.
type ExecOptions struct {
	InitDefaultDebugServer bool // enable default debug server.
	InitTracing            bool
	InitProfiler           bool
	FailOnValueError       bool

	LoadConfig    func(cmd *cobra.Command, vip *viper.Viper) error
	LoggerFactory func(*zap.Logger) *zap.Logger
}

// ExecWithCustomOptions runs a Cobra command with custom options.
func ExecWithCustomOptions(cmd *cobra.Command, opts ExecOptions) {
	cmd.AddCommand(&cobra.Command{
		Use:         "version",
		Short:       "output the version's build information, if any",
		RunE:        cmdVersion,
		Annotations: map[string]string{"type": "setup"}})

	exe, err := os.Executable()
	if err == nil && cmd.Use == "" {
		cmd.Use = exe
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	cleanup(cmd, &opts)
	err = cmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

// Ctx returns the appropriate context.Context for ExecuteWithConfig commands.
func Ctx(cmd *cobra.Command) (context.Context, context.CancelFunc) {
	commandMtx.Lock()
	defer commandMtx.Unlock()

	ctx := contexts[cmd]
	if ctx == nil {
		ctx = context.Background()
		contexts[cmd] = ctx
	}

	cancel := cancels[cmd]
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
		contexts[cmd] = ctx
		cancels[cmd] = cancel

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-c
			log.Printf("Got a signal from the OS: %q", sig)
			signal.Stop(c)
			cancel()
		}()
	}

	return ctx, cancel
}

// AtomicLevel returns the appropriate zap.AtomicLevel for ExecuteWithConfig commands.
func AtomicLevel(cmd *cobra.Command) *zap.AtomicLevel {
	commandMtx.Lock()
	defer commandMtx.Unlock()

	atomicLevel := atomicLevels[cmd]

	return atomicLevel
}

// Viper returns the appropriate *viper.Viper for the command, creating if necessary.
func Viper(cmd *cobra.Command) (*viper.Viper, error) {
	return ViperWithCustomConfig(cmd, LoadConfig)
}

// ViperWithCustomConfig returns the appropriate *viper.Viper for the command, creating if necessary. Custom
// config load logic can be defined with "loadConfig" parameter.
func ViperWithCustomConfig(cmd *cobra.Command, loadConfig func(cmd *cobra.Command, vip *viper.Viper) error) (*viper.Viper, error) {
	commandMtx.Lock()
	defer commandMtx.Unlock()

	if vip := vipers[cmd]; vip != nil {
		return vip, nil
	}

	vip := viper.New()
	if err := vip.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	prefix := os.Getenv("STORJ_ENV_PREFIX")
	if prefix == "" {
		prefix = "storj"
	}

	vip.SetEnvPrefix(prefix)
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	vip.AutomaticEnv()

	err := loadConfig(cmd, vip)
	if err != nil {
		return nil, err
	}

	vipers[cmd] = vip
	return vip, nil
}

// LoadConfig loads configuration into *viper.Viper from file specified with "config-dir" flag.
func LoadConfig(cmd *cobra.Command, vip *viper.Viper) error {
	cfgFlag := cmd.Flags().Lookup("config-dir")
	if cfgFlag != nil && cfgFlag.Value.String() != "" {
		path := filepath.Join(os.ExpandEnv(cfgFlag.Value.String()), DefaultCfgFilename)
		if fileExists(path) {
			setupCommand := cmd.Annotations["type"] == "setup"
			vip.SetConfigFile(path)
			if err := vip.ReadInConfig(); err != nil && !setupCommand {
				return err
			}
		}
	}
	return nil
}

var traceOut = flag.String("debug.trace-out", "", "If set, a path to write a process trace SVG to")

func getRoot(cmd *cobra.Command) *cobra.Command {
	if cmd.HasParent() {
		return getRoot(cmd.Parent())
	}
	return cmd
}

func cleanup(cmd *cobra.Command, opts *ExecOptions) {
	for _, ccmd := range cmd.Commands() {
		cleanup(ccmd, opts)
	}
	if cmd.Run != nil {
		panic("Please use cobra's RunE instead of Run")
	}
	internalRun := cmd.RunE
	if internalRun == nil {
		return
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		ctx := context.Background()
		defer mon.TaskNamed("root")(&ctx)(&err)

		vip, err := ViperWithCustomConfig(cmd, opts.LoadConfig)
		if err != nil {
			return err
		}

		commandMtx.Lock()
		configValues := configs[cmd]
		commandMtx.Unlock()

		var (
			brokenKeys  = map[string]struct{}{}
			missingKeys = map[string]struct{}{}
			usedKeys    = map[string]struct{}{}
			allKeys     = map[string]struct{}{}
			allSettings = vip.AllSettings()
		)

		for _, config := range configValues {
			// Decode and all of the resulting keys into our sets
			res := structs.Decode(allSettings, config)
			for key := range res.Used {
				usedKeys[key] = struct{}{}
				allKeys[key] = struct{}{}
			}
			for key := range res.Missing {
				missingKeys[key] = struct{}{}
				allKeys[key] = struct{}{}
			}
			for key := range res.Broken {
				brokenKeys[key] = struct{}{}
				allKeys[key] = struct{}{}
			}
		}

		// Propagate keys that are missing to flags, and remove any used keys
		// from the missing set.
		for key := range missingKeys {
			if f := cmd.Flags().Lookup(key); f != nil {
				val := vip.GetString(key)
				err := f.Value.Set(val)
				f.Changed = val != f.DefValue
				if err != nil {
					brokenKeys[key] = struct{}{}
				} else {
					usedKeys[key] = struct{}{}
				}
			} else if f := flag.Lookup(key); f != nil {
				err := f.Value.Set(vip.GetString(key))
				if err != nil {
					brokenKeys[key] = struct{}{}
				} else {
					usedKeys[key] = struct{}{}
				}
			}
		}
		for key := range missingKeys {
			if _, ok := usedKeys[key]; ok {
				delete(missingKeys, key)
			}
		}

		logger, atomicLevel, err := NewLogger(getRoot(cmd).Use)
		if err != nil {
			return err
		}
		if opts.LoggerFactory != nil {
			logger = opts.LoggerFactory(logger)
		}

		commandMtx.Lock()
		atomicLevels[cmd] = atomicLevel
		commandMtx.Unlock()

		if vip.ConfigFileUsed() != "" {
			path, err := filepath.Abs(vip.ConfigFileUsed())
			if err != nil {
				path = vip.ConfigFileUsed()
				logger.Debug("unable to resolve path", zap.Error(err))
			}

			logger.Info("Configuration loaded", zap.String("Location", path))
		}

		defer func() { _ = logger.Sync() }()
		defer zap.ReplaceGlobals(logger)()
		defer zap.RedirectStdLog(logger)()

		// okay now that logging is working, inform about the broken keys
		if cmd.Annotations["type"] != "helper" {
			for key := range missingKeys {
				logger.Info("Invalid configuration file key", zap.String("Key", key))
			}
		}
		for key := range brokenKeys {
			if opts.FailOnValueError {
				return errs.New("Invalid configuration file value for key: %s", key)
			}
			logger.Info("Invalid configuration file value for key", zap.String("Key", key))
		}

		if opts.InitTracing {
			var certPath string
			certPathFlag := cmd.Flags().Lookup("identity.cert-path")
			if certPathFlag != nil {
				certPath = certPathFlag.Value.String()
			}
			collector, unregister, err := InitTracingWithHostname(ctx, logger, nil, certPath)
			if err != nil {
				logger.Debug("failed to initialize tracing collector", zap.Error(err))
			} else if collector != nil {
				ctx, cancel := context.WithCancel(ctx)
				var eg errgroup.Group
				eg.Go(func() error {
					collector.Run(ctx)
					return nil
				})
				defer func() {
					cancel()
					unregister()
					err := eg.Wait()
					if err != nil {
						logger.Debug("failed to stop tracing collector", zap.Error(err))
					}
				}()
			}
		}

		if opts.InitDefaultDebugServer {
			err = initDebug(logger, monkit.Default, atomicLevel)
			if err != nil {
				withoutStack := errors.New(err.Error())
				logger.Debug("failed to start debug endpoints", zap.Error(withoutStack))
				err = nil
			}
		}

		if opts.InitProfiler && initProfiler != nil {
			profilerErr := initProfiler(logger)
			if profilerErr != nil {
				logger.Debug("failed to init debug profiler", zap.Error(profilerErr))
			}
		}

		var workErr error
		work := func(ctx context.Context) {
			commandMtx.Lock()
			contexts[cmd] = ctx
			commandMtx.Unlock()
			defer func() {
				commandMtx.Lock()
				delete(contexts, cmd)
				delete(cancels, cmd)
				commandMtx.Unlock()
			}()

			workErr = internalRun(cmd, args)
		}

		if *traceOut != "" {
			fh, err := os.Create(*traceOut)
			if err != nil {
				return err
			}
			if strings.HasSuffix(*traceOut, ".json") {
				err = present.SpansToJSON(fh, collect.CollectSpans(ctx, work))
			} else {
				err = present.SpansToSVG(fh, collect.CollectSpans(ctx, work))
			}
			err = errs.Combine(err, fh.Close())
			if err != nil {
				logger.Error("failed to write svg", zap.Error(err))
			}
		} else {
			work(ctx)
		}

		err = workErr
		if err != nil {
			// If our log output is stdout or stderr, we should log at the debug
			// level so that the stack trace is available for debugging, but
			// if debug is not enabled, print a nice error for the user to read.
			// If it's going anywhere else then we want to be very loud, so use Fatal.

			switch *logOutput {
			case "stdout":
				logger.Debug("Unrecoverable error", zap.Error(err))
				fmt.Fprintln(os.Stdout, "Error:", err.Error())
			case "stderr":
				logger.Debug("Unrecoverable error", zap.Error(err))
				fmt.Fprintln(os.Stderr, "Error:", err.Error())
			default:
				logger.Fatal("Unrecoverable error", zap.Error(err))
			}
			_ = logger.Sync()
			os.Exit(1)
		}

		return nil
	}
}

func cmdVersion(cmd *cobra.Command, args []string) (err error) {
	fmt.Println(version.Build)
	return nil
}
