package types

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
