package objects

import "rafal.dev/objects/simple"

var (
	ErrOutOfBounds    = simple.ErrOutOfBounds
	ErrNotFound       = simple.ErrNotFound
	ErrUnexpectedType = simple.ErrUnexpectedType
)

type (
	Error = simple.Error
)
