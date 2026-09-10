// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information

package sync2

import (
	"errors"
	"os"
	"sync/atomic"
)

// sharedFile implements Read, WriteAt offset to the file.
//
// The file is closed by the tee once every handle has been closed.
type sharedFile struct {
	file  *os.File
	read  int64
	write int64
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
func (buf *sharedFile) Close() error { return nil }

type memoryBlock struct {
	offset int
	data   []byte
	next   *memoryBlock
}

// blockReader implements io.ReadCloser on a memoryBlock.
//
// current is atomic because teeReader reads from the buffer without holding
// the tee lock, so Close can run concurrently with Read.
type blockReader struct {
	current atomic.Pointer[memoryBlock]
	read    int
}

func newBlockReader(block *memoryBlock) *blockReader {
	buf := &blockReader{}
	buf.current.Store(block)
	return buf
}

// blockWriter implements io.WriteCloser on a memoryBlock.
type blockWriter struct {
	current *memoryBlock
	write   int
}

var (
	errReaderPassedWriter = errors.New("block reader passed writer")
	errWriterMissingBlock = errors.New("block writer is missing a block")
)

// Reader implements io.Reader method.
func (buf *blockReader) Read(data []byte) (amount int, err error) {
	into := data
	for len(into) > 0 {
		cur := buf.current.Load()
		// if we don't have a block, we've finished the data
		if cur == nil {
			return amount, errReaderPassedWriter
		}

		// check whether we should proceed to the next block,
		// without undoing a concurrent Close
		if buf.read-cur.offset >= len(cur.data) {
			buf.current.CompareAndSwap(cur, cur.next)
			continue
		}

		// copy as much as we can
		n := copy(into, cur.data[buf.read-cur.offset:])

		into = into[n:]
		amount += n
		buf.read += n
	}
	return amount, nil
}

// Write implements io.Writer method.
func (buf *blockWriter) Write(data []byte) (amount int, err error) {
	out := data
	for len(out) > 0 {
		cur := buf.current
		// if we don't have a block, there's an error
		if cur == nil {
			return amount, errWriterMissingBlock
		}

		// we've reached end of the current block
		if buf.write-cur.offset >= len(cur.data) {
			cur.next = &memoryBlock{
				offset: cur.offset + len(cur.data),
				data:   make([]byte, len(cur.data)),
			}
			buf.current = cur.next
			continue
		}

		// copy as much as we can
		n := copy(cur.data[buf.write-cur.offset:], out)

		out = out[n:]
		amount += n
		buf.write += n
	}
	return amount, nil
}

// Close implements io.Closer methods.
//
// Dropping the block reference lets the blocks the reader did not consume be
// garbage collected.
func (buf *blockReader) Close() error {
	buf.current.Store(nil)
	return nil
}

// Close implements io.Closer methods.
func (buf *blockWriter) Close() error { return nil }
