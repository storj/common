// Code generated by protoc-gen-pico. DO NOT EDIT.
// source: encryption.proto
//
// versions:
//     protoc-gen-pico: v0.0.3
//     protoc:          v5.27.3

package pb

import (
	picobuf "storj.io/picobuf"
	strconv "strconv"
)

type CipherSuite int32

const (
	CipherSuite_ENC_UNSPECIFIED CipherSuite = 0
	CipherSuite_ENC_NULL        CipherSuite = 1
	CipherSuite_ENC_AESGCM      CipherSuite = 2
	CipherSuite_ENC_SECRETBOX   CipherSuite = 3
)

func (m CipherSuite) String() string {
	switch m {
	case CipherSuite_ENC_UNSPECIFIED:
		return "ENC_UNSPECIFIED"
	case CipherSuite_ENC_NULL:
		return "ENC_NULL"
	case CipherSuite_ENC_AESGCM:
		return "ENC_AESGCM"
	case CipherSuite_ENC_SECRETBOX:
		return "ENC_SECRETBOX"
	default:
		return "CipherSuite(" + strconv.Itoa(int(m)) + ")"
	}
}

type EncryptionParameters struct {
	CipherSuite CipherSuite `json:"cipher_suite,omitempty"`
	BlockSize   int64       `json:"block_size,omitempty"`
}

func (m *EncryptionParameters) Encode(c *picobuf.Encoder) bool {
	if m == nil {
		return false
	}
	c.Int32(1, (*int32)(&m.CipherSuite))
	c.Int64(2, &m.BlockSize)
	return true
}

func (m *EncryptionParameters) Decode(c *picobuf.Decoder) {
	if m == nil {
		return
	}
	c.Int32(1, (*int32)(&m.CipherSuite))
	c.Int64(2, &m.BlockSize)
}
