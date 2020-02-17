package json

import (
	r "github.com/freerware/tutor/api/representations"
)

type Error struct {
	r.Representation

	// Message represents a user friendly message.
	Message string `json:"message"`

	// DetailedMessage represents an engineer friendly message.
	DetailedMessage string `json:"detailedMessage"`
}

// Bytes provides the representation as bytes.
func (e Error) Bytes() ([]byte, error) {
	return e.Base.Bytes(&e)
}

// FromBytes constructs the representation from bytes.
func (e Error) FromBytes(b []byte) error {
	return e.Base.FromBytes(b, &e)
}
