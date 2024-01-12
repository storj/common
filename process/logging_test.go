// Copyright (C) 2024 Storj Labs, Inc.
// See LICENSE for copying information.

package process

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNamed(t *testing.T) {
	ptr := func(val string) *string {
		return &val
	}

	customLevel = ptr("log1=WARN")
	defer func() {
		customLevel = ptr("")
	}()

	out := bytes.NewBuffer([]byte{})
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(out), zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})))

	NamedLog(logger, "asd").Info("ahoj")
	require.True(t, strings.Contains(out.String(), "ahoj"))
	out.Reset()

	NamedLog(logger, "log1").Info("ahoj")
	require.False(t, strings.Contains(out.String(), "ahoj"))
	out.Reset()

	NamedLog(logger, "log1").Warn("ahoj")
	require.True(t, strings.Contains(out.String(), "ahoj"))
	out.Reset()

	customLevel = ptr("")
	NamedLog(logger, "log1").Warn("ahoj")
	require.False(t, strings.Contains(out.String(), "Invalid"))
	out.Reset()
}
