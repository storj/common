// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package macaroon

import (
	"crypto/rand"
	"encoding/json"

	"storj.io/common/encryption"
	"storj.io/common/storj"
)

// NewCaveat returns a Caveat with a random generated nonce.
func NewCaveat() (Caveat, error) {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	return Caveat{Nonce: buf[:]}, err
}

type caveatPathMarshal struct {
	Bucket              string `json:"bucket,omitempty"`
	EncryptedPathPrefix string `json:"encrypted_path_prefix,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (cp *Caveat_Path) MarshalJSON() ([]byte, error) {
	key, err := storj.NewKey([]byte{})
	if err != nil {
		return nil, err
	}

	prefix, err := encryption.DecryptPathRaw(string(cp.EncryptedPathPrefix), storj.EncNullBase64URL, key)
	if err != nil {
		return nil, err
	}

	return json.Marshal(caveatPathMarshal{
		Bucket:              string(cp.Bucket),
		EncryptedPathPrefix: prefix,
	})
}
