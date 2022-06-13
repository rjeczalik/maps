package objects

import "rafal.dev/objects/types"

var (
	TypeMap    = types.TypeMap
	TypeSlice  = types.TypeSlice
	TypeStruct = types.TypeStruct
)

type (
	Type      = types.Type
	Reader    = types.Reader
	Writer    = types.Writer
	Interface = types.Interface
	Meta      = types.Meta
	Iter      = types.Iter
)

type (
	Key            = types.Key
	PrefixedWriter = types.PrefixedWriter
	PrefixedReader = types.PrefixedReader
	Prefixed       = types.Prefixed
)
