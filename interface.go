package objects

import "rafal.dev/objects/types"

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
	MapType    = types.MapType
	SliceType  = types.SliceType
	StructType = types.StructType
)
