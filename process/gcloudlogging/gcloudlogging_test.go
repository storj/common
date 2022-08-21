// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package gcloudlogging

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestEncoder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	core := zapcore.NewCore(NewEncoder(NewEncoderConfig()), zapcore.AddSync(buf), zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	var got struct {
		HTTPRequest interface{}            `json:"httpRequest"`
		Labels      map[string]interface{} `json:"logging.googleapis.com/labels"`
		Message     string                 `json:"message"`
		Severity    string                 `json:"severity"`
		Time        time.Time              `json:"time"`
	}

	logger.Debug("a", zap.Bool("b", true))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))
	assert.Nil(t, got.HTTPRequest)
	assert.Nil(t, got.Labels["httpRequest"])
	assert.Equal(t, true, got.Labels["b"])
	assert.Nil(t, got.Labels["name"])
	assert.NotEmpty(t, got.Labels["caller"])
	assert.Equal(t, "a", got.Message)
	assert.Equal(t, "DEBUG", got.Severity)
	assert.NotZero(t, got.Time)

	logger.Info("c", zap.String("d", "e"))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))
	assert.Nil(t, got.HTTPRequest)
	assert.Nil(t, got.Labels["httpRequest"])
	assert.Equal(t, "e", got.Labels["d"])
	assert.Nil(t, got.Labels["name"])
	assert.NotEmpty(t, got.Labels["caller"])
	assert.Equal(t, "c", got.Message)
	assert.Equal(t, "INFO", got.Severity)
	assert.NotZero(t, got.Time)

	logger.Warn("f", zap.Bool("g", false))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))
	assert.Nil(t, got.HTTPRequest)
	assert.Nil(t, got.Labels["httpRequest"])
	assert.Equal(t, false, got.Labels["g"])
	assert.Nil(t, got.Labels["name"])
	assert.NotEmpty(t, got.Labels["caller"])
	assert.Equal(t, "f", got.Message)
	assert.Equal(t, "WARNING", got.Severity)
	assert.NotZero(t, got.Time)

	logger.Named("h").Error("i", zap.Bool("httpRequest", true))
	require.NoError(t, logger.Sync())

	require.NoError(t, json.NewDecoder(buf).Decode(&got))
	assert.Equal(t, true, got.HTTPRequest)
	assert.Nil(t, got.Labels["httpRequest"])
	assert.Equal(t, "h", got.Labels["name"])
	assert.NotEmpty(t, got.Labels["caller"])
	assert.True(t, strings.HasPrefix(got.Message, "i"))
	assert.NotEqual(t, "i", got.Message)
	assert.Equal(t, "ERROR", got.Severity)
	assert.NotZero(t, got.Time)
}
