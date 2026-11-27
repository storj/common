// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package encryption

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/common/memory"
	"storj.io/common/storj"
	"storj.io/common/testrand"
)

func TestCalcEncompassingBlocks(t *testing.T) {
	for _, example := range []struct {
		blockSize                              int
		offset, length, firstBlock, blockCount int64
	}{
		{2, 3, 4, 1, 3},
		{4, 0, 0, 0, 0},
		{4, 0, 1, 0, 1},
		{4, 0, 2, 0, 1},
		{4, 0, 3, 0, 1},
		{4, 0, 4, 0, 1},
		{4, 0, 5, 0, 2},
		{4, 0, 6, 0, 2},
		{4, 1, 0, 0, 0},
		{4, 1, 1, 0, 1},
		{4, 1, 2, 0, 1},
		{4, 1, 3, 0, 1},
		{4, 1, 4, 0, 2},
		{4, 1, 5, 0, 2},
		{4, 1, 6, 0, 2},
		{4, 2, 0, 0, 0},
		{4, 2, 1, 0, 1},
		{4, 2, 2, 0, 1},
		{4, 2, 3, 0, 2},
		{4, 2, 4, 0, 2},
		{4, 2, 5, 0, 2},
		{4, 2, 6, 0, 2},
		{4, 3, 0, 0, 0},
		{4, 3, 1, 0, 1},
		{4, 3, 2, 0, 2},
		{4, 3, 3, 0, 2},
		{4, 3, 4, 0, 2},
		{4, 3, 5, 0, 2},
		{4, 3, 6, 0, 3},
		{4, 4, 0, 1, 0},
		{4, 4, 1, 1, 1},
		{4, 4, 2, 1, 1},
		{4, 4, 3, 1, 1},
		{4, 4, 4, 1, 1},
		{4, 4, 5, 1, 2},
		{4, 4, 6, 1, 2},
	} {
		first, count := CalcEncompassingBlocks(
			example.offset, example.length, example.blockSize)
		if first != example.firstBlock || count != example.blockCount {
			t.Fatalf("invalid calculation for %#v: %v %v", example, first, count)
		}
	}
}

type nopTransformer struct {
	blockSize int
}

func NopTransformer(blockSize int) Transformer {
	return &nopTransformer{blockSize: blockSize}
}

func (t *nopTransformer) InBlockSize() int {
	return t.blockSize
}

func (t *nopTransformer) OutBlockSize() int {
	return t.blockSize
}

func (t *nopTransformer) Transform(out, in []byte, blockNum int64) (
	[]byte, error) {
	out = append(out, in...)
	return out, nil
}

func TestTransformer(t *testing.T) {
	transformer := NopTransformer(4 * 1024)
	data := testrand.BytesInt(transformer.InBlockSize() * 10)

	transformed := TransformReader(
		io.NopCloser(bytes.NewReader(data)),
		transformer, 0)
	data2, err := io.ReadAll(transformed)
	if assert.NoError(t, err) {
		assert.Equal(t, data, data2)
	}
}

func TestTransformWriter(t *testing.T) {
	for range 100 {
		tr := NopTransformer(64)
		data := testrand.BytesInt(tr.InBlockSize()*100 + testrand.Intn(tr.InBlockSize()))

		var buf bytes.Buffer
		trw := TransformWriterPadded(&buf, tr)

		for d := data; len(d) > 0; {
			n := len(d)
			if largest := 10 * tr.InBlockSize(); n > largest {
				n = 10 * tr.InBlockSize()
			}
			n = testrand.Intn(n + 1)

			_, err := trw.Write(d[:n])
			require.NoError(t, err)
			d = d[n:]
		}

		require.NoError(t, trw.Close())

		// ensure that after Close no more writes can happen
		_, err := trw.Write([]byte{1})
		require.Error(t, err)
		require.NoError(t, trw.Close())

		pad := PadReader(io.NopCloser(bytes.NewBuffer(data)), tr.InBlockSize())
		trr := TransformReader(pad, tr, 0)
		exp, err := io.ReadAll(trr)
		require.NoError(t, err)

		require.Equal(t, exp, buf.Bytes())
	}
}

func TestTransformerSize(t *testing.T) {
	for i, tt := range []struct {
		blockSize     int
		blocks        int
		expectedSize  int64
		unexpectedEOF bool
	}{
		{4, 10, 0, false},
		{4, 10, 3 * 10, false},
		{4, 10, 4*10 - 1, false},
		{4, 10, 4 * 10, false},
		{4, 10, 4*10 + 1, true},
		{4, 10, 4 * 11, true},
	} {
		errTag := fmt.Sprintf("Test case #%d", i)
		transformer := NopTransformer(tt.blockSize)
		data := testrand.BytesInt(transformer.InBlockSize() * tt.blocks)
		transformed := TransformReaderSize(
			io.NopCloser(bytes.NewReader(data)),
			transformer, 0, tt.expectedSize)
		data2, err := io.ReadAll(transformed)
		if tt.unexpectedEOF {
			assert.EqualError(t, err, io.ErrUnexpectedEOF.Error(), errTag)
		} else if assert.NoError(t, err, errTag) {
			assert.Equal(t, data, data2, errTag)
		}
	}
}

func BenchmarkTransformer(b *testing.B) {
	forAllCiphers(func(cipher storj.CipherSuite) {
		for _, dataSize := range []memory.Size{
			0,
			1,
			1 * memory.KiB,
			32 * memory.KiB,
			1 * memory.MiB,
			4 * memory.MiB,
		} {
			b.Run(fmt.Sprintf("%v/%v", cipher, dataSize), func(b *testing.B) {
				b.SetBytes(dataSize.Int64())

				parameters := storj.EncryptionParameters{CipherSuite: cipher, BlockSize: 1 * memory.KiB.Int32()}

				calculatedSize, err := CalcEncryptedSize(dataSize.Int64(), parameters)
				require.NoError(b, err)

				for i := 0; i < b.N; i++ {
					encrypter, err := NewEncrypter(parameters.CipherSuite, new(storj.Key), new(storj.Nonce), int(parameters.BlockSize))
					require.NoError(b, err)

					randReader := io.NopCloser(io.LimitReader(testrand.Reader(), dataSize.Int64()))
					reader := TransformReader(PadReader(randReader, encrypter.InBlockSize()), encrypter, 0)

					cipherData, err := io.ReadAll(reader)
					assert.NoError(b, err)
					assert.EqualValues(b, calculatedSize, len(cipherData))
				}
			})
		}
	})
}

func testBlockTransformers(t testing.TB, cipher storj.CipherSuite, blockSize int) (enc, dec Transformer) {
	var key storj.Key
	for i := range key {
		key[i] = byte(i)
	}
	// low bytes near overflow, so block numbers exercise the carry
	nonce := storj.Nonce{0: 1, 1: 0xfe, 2: 0xff, 3: 0xff}

	enc, err := NewEncrypter(cipher, &key, &nonce, blockSize)
	require.NoError(t, err)
	dec, err = NewDecrypter(cipher, &key, &nonce, blockSize)
	require.NoError(t, err)
	return enc, dec
}

func TestBlockTransformGolden(t *testing.T) {
	for cipher, expected := range map[storj.CipherSuite]string{
		storj.EncAESGCM:    "69ed90be78e440265434c5823184736e24b537e058d50ea50a5659ebf53538c6",
		storj.EncSecretBox: "fa9758befadfe85fb9e57e0db9b69a385a3522e58898a5af7db49826998c960d",
	} {
		enc, dec := testBlockTransformers(t, cipher, 64)
		plain := bytes.Repeat([]byte{0xa5}, enc.InBlockSize())

		h := sha256.New()
		for _, blockNum := range []int64{0, 1, 2, 255, 256, 1 << 40} {
			c, err := enc.Transform(nil, plain, blockNum)
			require.NoError(t, err)
			_, _ = h.Write(c)

			p, err := dec.Transform(nil, c, blockNum)
			require.NoError(t, err)
			require.Equal(t, plain, p)
		}
		assert.Equal(t, expected, hex.EncodeToString(h.Sum(nil)), cipher)
	}
}

func TestBlockTransformConcurrent(t *testing.T) {
	for _, cipher := range []storj.CipherSuite{storj.EncAESGCM, storj.EncSecretBox} {
		enc, dec := testBlockTransformers(t, cipher, 1024)
		plain := testrand.BytesInt(enc.InBlockSize())

		const n = 16
		expected := make([][]byte, n)
		for i := range expected {
			var err error
			expected[i], err = enc.Transform(nil, plain, int64(i))
			require.NoError(t, err)
		}

		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(blockNum int64) {
				defer wg.Done()
				for k := 0; k < 100; k++ {
					c, err := enc.Transform(nil, plain, blockNum)
					assert.NoError(t, err)
					assert.Equal(t, expected[blockNum], c)

					p, err := dec.Transform(nil, c, blockNum)
					assert.NoError(t, err)
					assert.Equal(t, plain, p)
				}
			}(int64(i))
		}
		wg.Wait()
	}
}

func BenchmarkBlockTransform(b *testing.B) {
	for _, cipher := range []storj.CipherSuite{storj.EncAESGCM, storj.EncSecretBox} {
		enc, dec := testBlockTransformers(b, cipher, 7424)
		plain := testrand.BytesInt(enc.InBlockSize())
		cipherData, err := enc.Transform(nil, plain, 0)
		require.NoError(b, err)

		b.Run(fmt.Sprintf("%v/encrypt", cipher), func(b *testing.B) {
			out := make([]byte, 0, enc.OutBlockSize())
			b.SetBytes(int64(len(plain)))
			b.ReportAllocs()
			for b.Loop() {
				_, _ = enc.Transform(out, plain, 0)
			}
		})
		b.Run(fmt.Sprintf("%v/decrypt", cipher), func(b *testing.B) {
			out := make([]byte, 0, dec.OutBlockSize())
			b.SetBytes(int64(len(cipherData)))
			b.ReportAllocs()
			for b.Loop() {
				_, _ = dec.Transform(out, cipherData, 0)
			}
		})
	}
}
