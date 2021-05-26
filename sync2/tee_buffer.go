// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package sync2

import (
	"io"
	"os"
	"sync/atomic"
)

// sharedFile implements Read, WriteAt offset to the file with reference counting.
type sharedFile struct {
	file  *os.File
	read  int64
	write int64
	open  *int64 // number of handles open
}

// ReadAt implements io.Reader methods.
func (buf *sharedFile) Read(data []byte) (amount int, err error) {
	amount, err = buf.file.ReadAt(data, buf.read)
	buf.read += int64(amount)
	return amount, err
}

// WriteAt implements io.Writer methods.
func (buf *sharedFile) Write(data []byte) (amount int, err error) {
	amount, err = buf.file.WriteAt(data, buf.write)
	buf.write += int64(amount)
	return amount, err
}

// Close implements io.Closer methods.
func (buf *sharedFile) Close() error {
	if atomic.AddInt64(buf.open, -1) == 0 {
		return buf.file.Close()
	}
	return nil
}

// sharedMemory implements Read, Write on a memory buffer.
type sharedMemory struct {
	memory []byte
	read   int
	write  int
}

// ReadAt implements io.Reader methods.
func (buf *sharedMemory) Read(data []byte) (amount int, err error) {
	if buf.read >= len(buf.memory) {
		return 0, io.ErrClosedPipe
	}
	amount = copy(data, buf.memory[buf.read:])
	buf.read += amount
	return amount, err
}

// WriteAt implements io.Writer methods.
func (buf *sharedMemory) Write(data []byte) (amount int, err error) {
	if buf.write >= len(buf.memory) {
		return 0, io.ErrClosedPipe
	}
	amount = copy(buf.memory[buf.write:], data)
	buf.write += amount
	if amount < len(data) {
		return amount, io.ErrShortWrite
	}
	return amount, err
}

// Close implements io.Closer methods.
func (buf *sharedMemory) Close() error { return nil }
