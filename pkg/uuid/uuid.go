package uuid

import "github.com/google/uuid"

type UUID = uuid.UUID

// New is xxx
func New() UUID {
	return uuid.New()
}

// NewMD5 is xxx
func NewMD5(space UUID, data []byte) UUID {
	return uuid.NewMD5(space, data)
}
