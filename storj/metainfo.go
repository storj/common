// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package storj

import "github.com/zeebo/errs"

var (
	// ErrNoPath is an error class for using empty path.
	ErrNoPath = errs.Class("no path specified")

	// ErrObjectNotFound is an error class for non-existing object.
	ErrObjectNotFound = errs.Class("object not found")
)

// ListDirection specifies listing direction.
type ListDirection int8

const (
	// Before lists backwards from cursor, without cursor [NOT SUPPORTED].
	Before = ListDirection(-2)
	// Backward lists backwards from cursor, including cursor [NOT SUPPORTED].
	Backward = ListDirection(-1)
	// Forward lists forwards from cursor, including cursor.
	Forward = ListDirection(1)
	// After lists forwards from cursor, without cursor.
	After = ListDirection(2)
)

// ListOptions lists objects.
type ListOptions struct {
	Prefix    Path
	Cursor    Path // Cursor is relative to Prefix, full path is Prefix + Cursor
	Delimiter rune
	Recursive bool
	Direction ListDirection
	Limit     int
	Status    int32
}
