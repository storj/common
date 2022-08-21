// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

// Package gcloudlogging provides special encoding, configuration for the
// encoder, and other constructs for go.uber.org/zap that make Cloud Logging
// understand its logs.
//
// All Storj-run applications will most certainly use the Cloud Logging agent
// instead of directly feeding Cloud Logging with LogEntries. This means we need
// to comply with the specification to make the message, level, time, and other
// fields gain special meaning that later allows us to construct powerful
// queries. Reference: https://cloud.google.com/logging/docs/structured-logging.
package gcloudlogging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type encoder struct {
	zapcore.Encoder
}

// EncodeEntry encodes entry and its fields, moving all fields except
// httpRequest (special case) as children of logging.googleapis.com/labels.
//
// Fields aren't top level and exist under 'labels' (see specification).
func (enc *encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var newFields []zapcore.Field

	for _, f := range fields {
		if f.Key == "httpRequest" {
			newFields = append(newFields, f)
			break
		}
	}

	newFields = append(newFields, zap.Object("logging.googleapis.com/labels", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		for _, f := range fields {
			if f.Key == "httpRequest" {
				// NOTE(artur): Object is marshaled lazily, so we can't add this
				// field to newFields here. Instead, we did it before.
				continue
			}
			f.AddTo(oe)
		}
		if ent.LoggerName != "" {
			zap.String("name", ent.LoggerName).AddTo(oe)
		} // It's better to have logger's name in labels.
		if c := ent.Caller.TrimmedPath(); c != "" {
			zap.String("caller", c).AddTo(oe)
		} // It's better to have caller in labels.
		return nil
	})))

	// Stack must be included in message. See:
	// https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
	if ent.Stack != "" {
		ent.Message += "\n" + ent.Stack
	}

	return enc.Encoder.EncodeEntry(ent, newFields)
}

// NewEncoder is like zapcore.NewJSONEncoder, but it moves fields and several
// keys in the log line so that Cloud Logging understands them better.
func NewEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &encoder{zapcore.NewJSONEncoder(cfg)}
}

// NewEncoderConfig creates zapcore.EncoderConfig suited for Cloud Logging.
func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "severity",
		TimeKey:        "time",
		NameKey:        "", // collapsed into labels
		CallerKey:      "", // collapsed into labels
		StacktraceKey:  "", // collapsed into messsage
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}
