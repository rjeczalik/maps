package objects

import (
	"reflect"

	"rafal.dev/objects/internal/misc"
	"rafal.dev/objects/types"
)

func Make(v any) Reader {
	if r, ok := v.(Reader); ok {
		return r
	}

	if v := types.Make(v); v != nil {
		return v
	}

	switch v := misc.ValueOf(v, true); v.Type().Kind() {
	case reflect.Struct:
		return &Struct{v: v}
	case reflect.Slice, reflect.Array:
		return &Slice{v: v}
	case reflect.Map:
		return &Map{v: v}
	}

	return nil
}

func tryMake(v any) any {
	if r := Make(v); r != nil {
		return r
	}
	return v
}
