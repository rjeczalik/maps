package mergelist

import "rafal.dev/objects/types"

type (
	Key   = types.Key
	Map   = types.Map
	Slice = types.Slice
)

type (
	Type       = types.Type
	Reader     = types.Reader
	SafeReader = types.SafeReader
	ListerTo   = types.ListerTo
	Writer     = types.Writer
	SafeWriter = types.SafeWriter
	Interface  = types.Interface
)

const (
	TypeMap    = types.TypeMap
	TypeSlice  = types.TypeSlice
	TypeStruct = types.TypeStruct
)

type all interface {
	Reader
	SafeReader
	ListerTo
	Writer
	SafeWriter
	Interface
}
