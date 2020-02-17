package yaml

import (
	r "github.com/freerware/tutor/api/representations"
)

type Error struct {
	r.Representation

	// Message represents a user friendly message.
	Message string `yaml:"message"`

	// DetailedMessage represents an engineer friendly message.
	DetailedMessage string `yaml:"detailedMessage"`
}

// Bytes provides the representation as bytes.
func (e Error) Bytes() ([]byte, error) {
	return e.Base.Bytes(&e)
}

// FromBytes constructs the representation from bytes.
func (e Error) FromBytes(b []byte) error {
	return e.Base.FromBytes(b, &e)
}
