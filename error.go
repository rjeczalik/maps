package objects

import "rafal.dev/objects/types"

var (
	ErrOutOfBounds    = types.ErrOutOfBounds
	ErrNotFound       = types.ErrNotFound
	ErrUnexpectedType = types.ErrUnexpectedType
)

type (
	Error = types.Error
)
