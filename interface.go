package objects

import "rafal.dev/objects/types"

type (
	Type          = types.Type
	Reader        = types.Reader
	SafeReader    = types.SafeReader
	ListerTo      = types.ListerTo
	Writer        = types.Writer
	SafeWriter    = types.SafeWriter
	Interface     = types.Interface
	SafeInterface = types.SafeInterface
	Iter          = types.Iter
)

const (
	TypeMap    = types.TypeMap
	TypeSlice  = types.TypeSlice
	TypeStruct = types.TypeStruct
)

type (
	Key            = types.Key
	PrefixedWriter = types.PrefixedWriter
	PrefixedReader = types.PrefixedReader
	Prefixed       = types.Prefixed
)
