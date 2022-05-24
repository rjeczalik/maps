package objects

import "rafal.dev/objects/simple"

type (
	Type       = simple.Type
	Reader     = simple.Reader
	SafeReader = simple.SafeReader
	ListerTo   = simple.ListerTo
	Writer     = simple.Writer
	SafeWriter = simple.SafeWriter
	Interface  = simple.Interface
)

const (
	MapType    = simple.MapType
	SliceType  = simple.SliceType
	StructType = simple.StructType
)
