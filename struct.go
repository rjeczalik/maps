package objects

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

var DefaultOptions = &Options{
	StructField: DefaultField,
}

type Options struct {
	StructField func(reflect.StructField) string
}

type Struct struct {
	v reflect.Value
}

var (
	_ Reader = (*Struct)(nil)
	_ Meta   = (*Struct)(nil)
)

func (s *Struct) Type() Type {
	return TypeStruct
}

func (s *Struct) Get(ctx context.Context, key string) (any, error) {
	switch v := s.v.FieldByName(key); {
	case !v.IsValid() || v.IsZero():
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: ErrNotFound,
		}
	case !v.CanInterface():
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Got: v,
			Err: fmt.Errorf("cannot access value: %s", v.Type()),
		}
	default:
		return tryMake(v.Interface()), nil
	}
}

func (s *Struct) List(ctx context.Context) ([]string, error) {
	var keys []string

	for _, f := range reflect.VisibleFields(s.v.Type()) {
		keys = append(keys, s.options().StructField(f))
	}

	return keys, nil
}

func (s *Struct) options() *Options {
	return DefaultOptions
}

func DefaultField(f reflect.StructField) string {
	return nonempty(
		tag(f.Tag, "object"),
		tag(f.Tag, "json"),
		tag(f.Tag, "yaml"),
		f.Name,
	)
}

func tag(t reflect.StructTag, name string) string {
	var s = t.Get(name)
	if i := strings.IndexRune(s, ','); i != -1 {
		s = s[:i]
	}
	return s
}

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}
