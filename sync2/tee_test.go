// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information

package sync2_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"storj.io/common/sync2"
	"storj.io/common/testrand"
)

func TestTee_Basic(t *testing.T) {
	run(t, func(t *testing.T, readers []sync2.PipeReader, writer sync2.PipeWriter) {
		var group errgroup.Group
		group.Go(func() error {
			n, err := writer.Write([]byte{1, 2, 3})
			assert.Equal(t, n, 3)
			assert.NoError(t, err)

			n, err = writer.Write([]byte{1, 2, 3})
			assert.Equal(t, n, 3)
			assert.NoError(t, err)

			assert.NoError(t, writer.Close())
			return nil
		})

		for i := 0; i < len(readers); i++ {
			i := i
			group.Go(func() error {
				data, err := ioutil.ReadAll(readers[i])
				assert.Equal(t, []byte{1, 2, 3, 1, 2, 3}, data)
				if err != nil {
					assert.Equal(t, io.EOF, err)
				}
				assert.NoError(t, readers[i].Close())
				return nil
			})
		}

		assert.NoError(t, group.Wait())
	})
}

func TestTee_CloseWithError(t *testing.T) {
	run(t, func(t *testing.T, readers []sync2.PipeReader, writer sync2.PipeWriter) {
		var failure = errors.New("write failure")

		var group errgroup.Group
		group.Go(func() error {
			n, err := writer.Write([]byte{1, 2, 3})
			assert.Equal(t, n, 3)
			assert.NoError(t, err)

			err = writer.CloseWithError(failure)
			assert.NoError(t, err)

			return nil
		})

		for i := 0; i < len(readers); i++ {
			i := i
			group.Go(func() error {
				_, err := ioutil.ReadAll(readers[i])
				if err != nil {
					assert.Equal(t, failure, err)
				}
				assert.NoError(t, readers[i].Close())
				return nil
			})
		}

		assert.NoError(t, group.Wait())
	})
}

const testBlockSize = 1024 // 1KiB

func TestTee_Blocks(t *testing.T) {
	run(t, func(t *testing.T, readers []sync2.PipeReader, writer sync2.PipeWriter) {
		// a single block
		expected1 := testrand.Bytes(testBlockSize)

		{
			n, err := writer.Write(expected1)
			require.NoError(t, err)
			require.Equal(t, testBlockSize, n)
		}

		eqint := func(a, b int) error {
			if a != b {
				return errs.New("values different %v != %v", a, b)
			}
			return nil
		}

		eqdata := func(a, b []byte) error {
			if !bytes.Equal(a, b) {
				return errs.New("data different")
			}
			return nil
		}

		concurrent(t,
			func() error {
				// read a single block
				full := make([]byte, testBlockSize)
				n, err := readers[0].Read(full)
				return errs.Combine(err,
					eqint(testBlockSize, n),
					eqdata(expected1, full),
				)
			}, func() error {
				// read a in half blocks
				half := make([]byte, testBlockSize/2)
				for k := 0; k < len(expected1); k += len(half) {
					n, err := readers[1].Read(half)
					err = errs.Combine(err,
						eqint(len(half), n),
						eqdata(expected1[k:k+len(half)], half),
					)
					if err != nil {
						return err
					}
				}
				return nil
			},
		)

		// multiple blocks
		expected2 := testrand.Bytes(testBlockSize * 2)

		{ // write two blocks
			n, err := writer.Write(expected2)
			require.NoError(t, err)
			require.Equal(t, 2048, n)
		}

		concurrent(t,
			func() error {
				// read two block
				full := make([]byte, 2*testBlockSize)
				n, err := readers[0].Read(full)
				return errs.Combine(err,
					eqint(testBlockSize*2, n),
					eqdata(expected2, full),
				)
			}, func() error {
				// read a in half blocks
				half := make([]byte, testBlockSize/2)
				for k := 0; k < len(expected2); k += len(half) {
					n, err := readers[1].Read(half)
					err = errs.Combine(err,
						eqint(len(half), n),
						eqdata(expected2[k:k+len(half)], half),
					)
					if err != nil {
						return err
					}
				}
				return nil
			},
		)

		// concurrent read write 1 byte
		concurrent(t,
			func() error {
				// write one byte
				n, err := writer.Write([]byte{1})
				return errs.Combine(err,
					eqint(1, n),
				)
			},
			func() error {
				// read across block boundary
				abyte := make([]byte, 1)
				n, err := readers[0].Read(abyte)
				return errs.Combine(err,
					eqint(1, n),
					eqdata([]byte{1}, abyte),
				)
			}, func() error {
				// read across block boundary
				abyte := make([]byte, 1)
				n, err := readers[1].Read(abyte)
				return errs.Combine(err,
					eqint(1, n),
					eqdata([]byte{1}, abyte),
				)
			},
		)

		// concurrent read&write multiple blocks
		expected3 := testrand.Bytes(testBlockSize * 3)
		concurrent(t,
			func() error {
				// write multiple blocks
				n, err := writer.Write(expected3)
				return errs.Combine(err,
					eqint(len(expected3), n),
				)
			},
			func() error {
				full := make([]byte, len(expected3))
				n, err := readers[0].Read(full)
				return errs.Combine(err,
					eqint(len(expected3), n),
					eqdata(expected3, full),
				)
			}, func() error {
				// read across block boundary
				full := make([]byte, len(expected3))
				n, err := readers[1].Read(full)
				return errs.Combine(err,
					eqint(len(expected3), n),
					eqdata(expected3, full),
				)
			},
		)
	})
}

func concurrent(t *testing.T, fns ...func() error) {
	var group errgroup.Group
	for _, fn := range fns {
		group.Go(fn)
	}
	require.NoError(t, group.Wait())
}

func run(t *testing.T, test func(t *testing.T, readers []sync2.PipeReader, writer sync2.PipeWriter)) {
	t.Parallel()

	t.Run("File", func(t *testing.T) {
		t.Parallel()

		readers, writer, err := sync2.NewTeeFile(2, "")
		if err != nil {
			t.Fatal(err)
		}

		test(t, readers, writer)
	})

	t.Run("Inmemory", func(t *testing.T) {
		t.Parallel()

		readers, writer, err := sync2.NewTeeInmemory(2, testBlockSize)
		if err != nil {
			t.Fatal(err)
		}

		test(t, readers, writer)
	})
}
