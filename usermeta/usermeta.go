// Copyright (C) 2025 Storj Labs, Inc.
// See LICENSE for copying information.

package usermeta

import (
	"encoding/json"
	"strings"

	"storj.io/common/pb"
)

// UserMeta stores simple string key-value metadata for objects.
type UserMeta map[string]string

// Marshal user metadata payload using the SerializableMeta structure.
func Marshal(meta UserMeta) ([]byte, error) {
	m := &pb.SerializableMeta{
		UserDefined: meta,
	}

	return pb.Marshal(m)
}

// Marshal a JSON-encoded, deeply structured user metadata using the
// SerializableMeta structure.
func MarshalJSON(meta string) ([]byte, error) {
	var dm DeepUserMeta
	err := json.Unmarshal([]byte(meta), &dm)
	if err != nil {
		return nil, err
	}

	m, err := dm.toUserMeta()
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

// UnmarshalJSON unmarshals the user metadata payload, converts it to a deeply
// structured metadata and returns it as a JSON string.
func UnmarshalJSON(data []byte) (string, error) {
	m, err := Unmarshal(data)
	if err != nil {
		return "", err
	}

	dm, err := m.toDeepUserMeta()
	if err != nil {
		return "", err
	}

	j, err := json.Marshal(dm)
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

// DeepUserMeta stores an arbitrary deep metadata structures for
// objects. It can be converted from/to one-level UserMeta objects.
type DeepUserMeta map[string]interface{}

func (m DeepUserMeta) toUserMeta() (UserMeta, error) {
	meta := make(UserMeta)

	for k, v := range m {
		if s, ok := v.(string); ok {
			meta[k] = s
		} else {
			j, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			meta["json:"+k] = string(j)
		}
	}
	return meta, nil
}

func (m UserMeta) toDeepUserMeta() (DeepUserMeta, error) {
	meta := make(DeepUserMeta)

	for k, v := range m {
		if strings.HasPrefix(k, "json:") {
			var i interface{}
			err := json.Unmarshal([]byte(v), &i)
			if err != nil {
				return nil, err
			}
			meta[k[5:]] = i
		} else {
			meta[k] = v
		}
	}
	return meta, nil
}
