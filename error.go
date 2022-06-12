package objects

import "rafal.dev/objects/types"

var (
	ErrOutOfBounds    = types.ErrOutOfBounds
	ErrEmpty          = types.ErrEmpty
	ErrNotFound       = types.ErrNotFound
	ErrNotDone        = types.ErrNotDone
	ErrUnexpectedType = types.ErrUnexpectedType
)

type (
	Error = types.Error
)
