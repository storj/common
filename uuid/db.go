// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package uuid

import (
	"database/sql/driver"
	"encoding/base64"
)

// Value implements sql/driver.Valuer interface.
func (uuid UUID) Value() (driver.Value, error) {
	return uuid[:], nil
}

// Scan implements sql.Scanner interface.
func (uuid *UUID) Scan(value interface{}) error {
	switch value := value.(type) {
	case []byte:
		x, err := FromBytes(value)
		if err != nil {
			return Error.Wrap(err)
		}
		*uuid = x
		return nil
	case string:
		x, err := FromString(value)
		if err != nil {
			return Error.Wrap(err)
		}
		*uuid = x
		return nil
	default:
		return Error.New("unable to scan %T into UUID", value)
	}
}

// DecodeSpanner implements spanner.Decoder.
func (uuid *UUID) DecodeSpanner(val any) (err error) {
	if v, ok := val.(string); ok {
		var buffer [16]byte
		b, err := base64.StdEncoding.AppendDecode(buffer[:0], []byte(v))
		if err != nil {
			return err
		}
		x, err := FromBytes(b)
		if err != nil {
			return Error.Wrap(err)
		}
		*uuid = x
		return nil
	}
	return uuid.Scan(val)
}

// EncodeSpanner implements spanner.Encoder.
func (uuid UUID) EncodeSpanner() (any, error) {
	return uuid.Value()
}

// NullUUID represents a UUID that may be null.
// NullUUID implements the Scanner interface so it can be used
// as a scan destination, similar to sql.NullString.
type NullUUID struct {
	UUID  UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Value implements sql/driver.Valuer interface.
func (n NullUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.UUID.Value()
}

// Scan implements sql.Scanner interface.
func (n *NullUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID, n.Valid = UUID{}, false
		return nil
	}

	// a NULL BYTES value gets returned from Spanner as an empty []byte
	if v, ok := value.([]byte); ok {
		if v == nil {
			n.UUID, n.Valid = UUID{}, false
			return nil
		}
	}

	n.Valid = true
	return n.UUID.Scan(value)
}

// EncodeSpanner implements spanner.Encoder.
func (n NullUUID) EncodeSpanner() (any, error) {
	return n.Value()
}

// DecodeSpanner implements spanner.Decoder.
func (n *NullUUID) DecodeSpanner(val any) (err error) {
	if v, ok := val.(string); ok {
		val, err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}
	}
	return n.Scan(val)
}
