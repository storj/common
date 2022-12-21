// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcloudlogging_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"storj.io/private/process/gcloudlogging"
)

func TestOperationEncoding(t *testing.T) {
	op := &gcloudlogging.Operation{
		ID:       "foo",
		Producer: "github.com/storj/gateway-mt",
		First:    true,
		Last:     false,
	}

	enc := gcloudlogging.NewEncoder(gcloudlogging.NewEncoderConfig())

	buf, err := enc.EncodeEntry(zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Message: "test",
		Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}, []zapcore.Field{
		gcloudlogging.LogOperation(op),
		zap.String("something", "else"),
	})
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"logging.googleapis.com/operation": {
			"first": true,
			"id": "foo",
			"producer": "github.com/storj/gateway-mt"
		},
		"logging.googleapis.com/labels": {
			"something": "else"
		},
		"logging.googleapis.com/severity": "INFO",
		"message": "test",
		"timestamp": {
			"nanos": 0,
			"seconds": 1.257894e+09
		}
	}`, buf.String())
}
