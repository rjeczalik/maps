package types

type Type string

const (
	MapType    Type = "Map"
	SliceType  Type = "Slice"
	StructType Type = "Struct"
)

func Make(v any) Interface {
	switch v := v.(type) {
	case Interface:
		return v
	case map[string]any:
		return Map(v)
	case []any:
		s := Slice(v)
		return &s
	case Map:
		return v
	case Slice:
		return &v
	case *Slice:
		return v
	default:
		return nil
	}
}

func tryMake(v any) any {
	if v := Make(v); v != nil {
		return v
	}
	return v
}

func makeOr(hint Type, def Interface) Interface {
	switch hint {
	case MapType, StructType:
		return make(Map)
	case SliceType:
		return &Slice{}
	default:
		return def
	}
}
