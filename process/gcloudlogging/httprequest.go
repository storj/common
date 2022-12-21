// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcloudlogging

import (
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPRequest represents HttpRequest field. See:
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
type HTTPRequest struct {
	RequestMethod                  string
	RequestURL                     string
	RequestSize                    int64
	Status                         int
	ResponseSize                   int64
	UserAgent                      string
	RemoteIP                       string
	ServerIP                       string
	Referer                        string
	Latency                        time.Duration
	CacheLookup                    bool
	CacheHit                       bool
	CacheValidatedWithOriginServer bool
	CacheFillBytes                 int64
	Protocol                       string
}

// MarshalLogObject implements zapcore.ObjectMarshaler.
// All fields are optional.
func (req *HTTPRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if req.RequestMethod != "" {
		enc.AddString("requestMethod", req.RequestMethod)
	}
	if req.RequestURL != "" {
		enc.AddString("requestUrl", req.RequestURL)
	}
	// note: RequestSize is header + body size in bytes, not just the body.
	if req.RequestSize != 0 {
		enc.AddString("requestSize", strconv.FormatInt(req.RequestSize, 10))
	}
	if req.Status != 0 {
		enc.AddInt("status", req.Status)
	}
	if req.ResponseSize != 0 {
		enc.AddString("responseSize", strconv.FormatInt(req.ResponseSize, 10))
	}
	if req.UserAgent != "" {
		enc.AddString("userAgent", req.UserAgent)
	}
	if req.RemoteIP != "" {
		enc.AddString("remoteIp", req.RemoteIP)
	}
	if req.ServerIP != "" {
		enc.AddString("serverIp", req.ServerIP)
	}
	if req.Referer != "" {
		enc.AddString("referer", req.Referer)
	}
	// spec wants it in seconds, e.g. "3.5s" with up to nine fractional digits.
	if req.Latency.Seconds() != 0 {
		enc.AddString("latency", strconv.FormatFloat(req.Latency.Seconds(), 'f', 9, 64)+"s")
	}

	// note: GCP logging treats bool false as empty.
	if req.CacheLookup {
		enc.AddBool("cacheLookup", true)
	}
	if req.CacheHit {
		enc.AddBool("cacheHit", true)
	}
	if req.CacheValidatedWithOriginServer {
		enc.AddBool("cacheValidatedWithOriginServer", true)
	}
	if req.CacheFillBytes != 0 {
		enc.AddString("cacheFillBytes", strconv.FormatInt(req.CacheFillBytes, 10))
	}
	if req.Protocol != "" {
		enc.AddString("protocol", req.Protocol)
	}

	return nil
}

// LogHTTPRequest returns a zapcore.Field for HTTPRequest.
func LogHTTPRequest(req *HTTPRequest) zapcore.Field {
	return zap.Object("logging.googleapis.com/httpRequest", req)
}
