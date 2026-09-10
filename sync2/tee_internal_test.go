// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information

package sync2

import (
	"io"
	"runtime"
	"testing"
	"weak"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// TestTeeInmemory_ClosedReaderReleasesBlocks checks that a reader which stops
// early stops pinning the rest of the stream. The in-memory tee keeps the
// blocks between the slowest reader and the writer alive, and the writer runs
// ahead of the furthest read request, so a reader that is closed part way
// through used to hold on to every block produced after it.
func TestTeeInmemory_ClosedReaderReleasesBlocks(t *testing.T) {
	const blockSize = 1024

	readers, writer, err := NewTeeInmemory(2, blockSize)
	require.NoError(t, err)

	early := readers[0].(*teeReader)
	buffer := early.buffer.(*blockReader)

	first := weak.Make(buffer.current.Load())
	require.NotNil(t, first.Value())

	// the writer only runs ahead of the furthest read request, so it has to be
	// running before anyone reads.
	var group errgroup.Group
	group.Go(func() error {
		defer func() { _ = writer.Close() }()
		_, err := writer.Write(make([]byte, blockSize*4))
		return err
	})

	// read a single byte, so the reader is left sitting on the first block.
	var one [1]byte
	_, err = io.ReadFull(early, one[:])
	require.NoError(t, err)
	require.NoError(t, early.Close())

	// drain the stream through the other reader, which drags the writer past
	// several blocks. the closed reader must not keep those blocks alive.
	_, err = io.ReadAll(readers[1])
	require.NoError(t, err)
	require.NoError(t, group.Wait())
	require.NoError(t, readers[1].Close())

	require.Nil(t, buffer.current.Load(), "closed reader still references a block")

	runtime.GC()
	require.Nil(t, first.Value(), "blocks held by the closed reader were not released")
	// the tee has to stay reachable, otherwise the block above is collected
	// whether or not Close released it.
	runtime.KeepAlive(buffer)
}
