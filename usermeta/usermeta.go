// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information.

package usermeta

import (
	"encoding/json"

	"storj.io/common/pb"
)

type UserMeta map[string]string

// Marshal user metadata payload using the SerializableMeta structure.
func Marshal(meta UserMeta) ([]byte, error) {
	m := &pb.SerializableMeta{
		UserDefined: meta,
	}

	return pb.Marshal(m)
}

// Marshal a JSON-encoded user metadata using the SerializableMeta structure.
func MarshalJSON(meta string) ([]byte, error) {
	var m UserMeta
	err := json.Unmarshal([]byte(meta), &m)
	if err != nil {
		return nil, err
	}

	return Marshal(m)
}

// Unmarshal user metadata payload using the SerializableMeta structure.
func Unmarshal(data []byte) (UserMeta, error) {
	m := new(pb.SerializableMeta)
	err := pb.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}

	return m.UserDefined, nil
}

// UnmarshalJSON unmarshals the user metadata payload and returns it as a JSON string.
func UnmarshalJSON(data []byte) (string, error) {
	meta, err := Unmarshal(data)
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

// Valid checks if the data is a valid user metadata payload.
func Valid(data []byte) bool {
	_, err := Unmarshal(data)
	return err == nil
}
