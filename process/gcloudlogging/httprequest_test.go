// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcloudlogging_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"storj.io/private/process/gcloudlogging"
)

func TestHTTPRequestEncoding(t *testing.T) {
	req := &gcloudlogging.HTTPRequest{
		RequestMethod: http.MethodGet,
		RequestURL:    "/lol",
		RequestSize:   12345,
		Status:        http.StatusBadGateway,
		ResponseSize:  0,
		UserAgent:     "lol/bot",
		RemoteIP:      "1.2.3.4",
		ServerIP:      "4.5.6.7",
		Referer:       "something",
		Latency:       1234567890123456789 * time.Nanosecond,
		Protocol:      "HTTP/1.1",
	}

	enc := gcloudlogging.NewEncoder(gcloudlogging.NewEncoderConfig())

	buf, err := enc.EncodeEntry(zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Message: "test",
		Time:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}, []zapcore.Field{
		gcloudlogging.LogHTTPRequest(req),
		zap.String("something", "else"),
	})
	require.NoError(t, err)

	assert.JSONEq(t, `{
		"logging.googleapis.com/httpRequest": {
			"latency": "1234567890.123456717s",
			"protocol": "HTTP/1.1",
			"referer": "something",
			"remoteIp": "1.2.3.4",
			"requestMethod": "GET",
			"requestSize": "12345",
			"requestUrl": "/lol",
			"serverIp": "4.5.6.7",
			"status": 502,
			"userAgent": "lol/bot"
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
