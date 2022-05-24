package objects

import (
	"fmt"
	"reflect"
	"strings"
)

type Struct struct {
	v reflect.Value
}

var (
	_ Reader     = (*Struct)(nil)
	_ SafeReader = (*Struct)(nil)
	_ ListerTo   = (*Struct)(nil)
)

func (s *Struct) Type() Type {
	return StructType
}

func (s *Struct) Get(key string) (any, bool) {
	v, err := s.SafeGet(key)
	return v, err == nil
}

func (s *Struct) SafeGet(key string) (any, error) {
	switch v := s.v.FieldByName(key); {
	case v.IsZero():
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: ErrNotFound,
		}
	case !v.CanInterface():
		fmt.Println("XD", key)
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

func (s *Struct) List() []string {
	var keys []string
	s.ListTo(&keys)
	return keys
}

func (s *Struct) ListTo(keys *[]string) {
	for _, f := range reflect.VisibleFields(s.v.Type()) {
		*keys = append(*keys, f.Name)
	}
}

func field(f reflect.StructField) string {
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
