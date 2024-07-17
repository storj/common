// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package gcloudlogging

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestEncoder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	core := zapcore.NewCore(NewEncoder(NewEncoderConfig()), zapcore.AddSync(buf), zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel)).Named("test")

	got := make(map[string]interface{})

	logger.Debug("a", zap.Bool("b", true))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	b, ok := got["b"].(bool)
	require.True(t, ok)
	assert.True(t, b)

	labels, ok := got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 1)
	assert.Equal(t, "test", labels["name"])

	timestamp, ok := got["timestamp"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, timestamp, 2)
	assert.NotZero(t, timestamp["seconds"])
	assert.NotZero(t, timestamp["nanos"])

	sourceLocation, ok := got["logging.googleapis.com/sourceLocation"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, sourceLocation, 3)
	assert.NotEmpty(t, sourceLocation["file"])
	assert.NotEmpty(t, sourceLocation["line"])
	assert.NotEmpty(t, sourceLocation["function"])

	assert.Equal(t, "a", got["message"])
	assert.Equal(t, "DEBUG", got["severity"])

	logger.Info("c", zap.String("d", "e"))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "e", labels["d"])

	assert.Equal(t, "c", got["message"])
	assert.Equal(t, "INFO", got["severity"])

	logger.Warn("f", zap.Stringer("g", bytes.NewBufferString("h")))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "h", labels["g"])

	assert.Equal(t, "f", got["message"])
	assert.Equal(t, "WARNING", got["severity"])

	logger.Error("i", zap.Error(errors.New("j")))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "j", labels["error"])

	assert.Equal(t, "i", got["message"])
	assert.Equal(t, "ERROR", got["severity"])
}

func TestEncoderChildLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	core := zapcore.NewCore(NewEncoder(NewEncoderConfig()), zapcore.AddSync(buf), zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel)).Named("test")
	childLogger := logger.With(zap.String("name", "child"))

	got := make(map[string]interface{})

	childLogger.Debug("a", zap.Bool("b", true))
	require.NoError(t, childLogger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	b, ok := got["b"].(bool)
	require.True(t, ok)
	assert.True(t, b)

	labels, ok := got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 1)
	assert.Equal(t, "test", labels["name"])

	timestamp, ok := got["timestamp"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, timestamp, 2)
	assert.NotZero(t, timestamp["seconds"])
	assert.NotZero(t, timestamp["nanos"])

	sourceLocation, ok := got["logging.googleapis.com/sourceLocation"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, sourceLocation, 3)
	assert.NotEmpty(t, sourceLocation["file"])
	assert.NotEmpty(t, sourceLocation["line"])
	assert.NotEmpty(t, sourceLocation["function"])

	assert.Equal(t, "a", got["message"])
	assert.Equal(t, "DEBUG", got["severity"])

	childLogger.Info("c", zap.String("d", "e"))
	require.NoError(t, childLogger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "e", labels["d"])

	assert.Equal(t, "c", got["message"])
	assert.Equal(t, "INFO", got["severity"])

	childLogger.Warn("f", zap.Stringer("g", bytes.NewBufferString("h")))
	require.NoError(t, childLogger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "h", labels["g"])

	assert.Equal(t, "f", got["message"])
	assert.Equal(t, "WARNING", got["severity"])

	childLogger.Error("i", zap.Error(errors.New("j")))
	require.NoError(t, childLogger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))

	labels, ok = got["logging.googleapis.com/labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Len(t, labels, 2)
	assert.Equal(t, "j", labels["error"])

	assert.Equal(t, "i", got["message"])
	assert.Equal(t, "ERROR", got["severity"])
}
