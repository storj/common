// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"

	"storj.io/private/cfgstruct"
)

var (
	// Error is a process error class.
	Error = errs.Class("process error")

	logLevel = zap.LevelFlag("log.level", func() zapcore.Level {
		if isDev() {
			return zapcore.DebugLevel
		}
		return zapcore.InfoLevel
	}(), "the minimum log level to log")
	logDev      = flag.Bool("log.development", isDev(), "if true, set logging to development mode")
	logCaller   = flag.Bool("log.caller", isDev(), "if true, log function filename and line number")
	logStack    = flag.Bool("log.stack", isDev(), "if true, log stack traces")
	logEncoding = flag.String("log.encoding", "", "configures log encoding. can either be 'console', 'json', or 'pretty'.")
	logOutput   = flag.String("log.output", "stderr", "can be stdout, stderr, or a filename")

	defaultLogEncoding = map[string]string{
		"uplink": "pretty",
	}
)

func init() {
	winFileSink := func(u *url.URL) (zap.Sink, error) {
		// Remove leading slash left by url.Parse()
		return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	}
	err := zap.RegisterSink("winfile", winFileSink)
	if err != nil {
		panic("Unable to register winfile sink: " + err.Error())
	}

	err = zap.RegisterEncoder("pretty", func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newPrettyEncoder(config), nil
	})
	if err != nil {
		panic("Unable to register pretty encoder: " + err.Error())
	}
}

func isDev() bool { return cfgstruct.DefaultsType() != "release" }

// NewLogger creates new logger configured by the process flags.
func NewLogger(processName string) (*zap.Logger, *zap.AtomicLevel, error) {
	return NewLoggerWithOutputPathsAndAtomicLevel(processName, *logOutput)
}

// NewLoggerWithOutputPaths is the same as NewLogger, but overrides the log output paths.
func NewLoggerWithOutputPaths(processName string, outputPaths ...string) (*zap.Logger, error) {
	logger, _, err := NewLoggerWithOutputPathsAndAtomicLevel(processName, outputPaths...)
	return logger, err
}

// NewLoggerWithOutputPathsAndAtomicLevel is the same as NewLoggerWithOutputPaths, but overrides the log output paths
// and returns the AtomicLevel.
func NewLoggerWithOutputPathsAndAtomicLevel(processName string, outputPaths ...string) (*zap.Logger, *zap.AtomicLevel, error) {
	levelEncoder := zapcore.CapitalLevelEncoder
	timeKey := "T"
	if os.Getenv("STORJ_LOG_NOTIME") != "" {
		// using environment variable STORJ_LOG_NOTIME to avoid additional flags
		timeKey = ""
	}

	encoding := *logEncoding
	if encoding == "" {
		encoding = defaultLogEncoding[processName]
		if encoding == "" {
			encoding = "console"
		}
	}

	atomicLevel := zap.NewAtomicLevelAt(*logLevel)
	logger, err := zap.Config{
		Level:             atomicLevel,
		Development:       *logDev,
		DisableCaller:     !*logCaller,
		DisableStacktrace: !*logStack,
		Encoding:          encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        timeKey,
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    levelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: outputPaths,
	}.Build()

	return logger, &atomicLevel, err
}

type prettyEncoder struct {
	*zapcore.MapObjectEncoder
	config zapcore.EncoderConfig
	pool   buffer.Pool
}

func newPrettyEncoder(config zapcore.EncoderConfig) *prettyEncoder {
	return &prettyEncoder{
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
		config:           config,
		pool:             buffer.NewPool(),
	}
}

func (p *prettyEncoder) Clone() zapcore.Encoder {
	rv := newPrettyEncoder(p.config)
	for key, val := range p.MapObjectEncoder.Fields {
		rv.MapObjectEncoder.Fields[key] = val
	}
	return rv
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (p *prettyEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	b := p.pool.Get()

	fmt.Fprintf(b, "%s\t%s\t%s\n",
		entry.Time.Format("15:04:05.000"),
		levelDecorate(entry.Level, entry.Level.CapitalString()),
		entry.Message)

	for _, field := range fields {
		m := zapcore.NewMapObjectEncoder()
		field.AddTo(m)
		for _, key := range sortedKeys(m.Fields) {
			if key == "errorVerbose" && !*logDev {
				continue
			}
			fmt.Fprintf(b, "\t%s: %s\n",
				key,
				strings.ReplaceAll(fmt.Sprint(m.Fields[key]), "\n", "\n\t"))
		}
	}

	return b, nil
}
