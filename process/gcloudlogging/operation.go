// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

//zapfields:ignore-file
package gcloudlogging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Operation represents LogEntryOperation field. See:
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
type Operation struct {
	ID       string
	Producer string
	First    bool
	Last     bool
}

// MarshalLogObject implements zapcore.ObjectMarshaler.
// All fields are optional.
func (op *Operation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if op.ID != "" {
		enc.AddString("id", op.ID)
	}
	if op.Producer != "" {
		enc.AddString("producer", op.Producer)
	}
	// note: GCP logging treats bool false as empty.
	if op.First {
		enc.AddBool("first", true)
	}
	if op.Last {
		enc.AddBool("last", true)
	}

	return nil
}

// LogOperation returns a zapcore.field for Operation.
func LogOperation(op *Operation) zapcore.Field {
	return zap.Object("logging.googleapis.com/operation", op)
}
