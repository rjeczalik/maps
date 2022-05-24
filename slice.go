package objects

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Slice struct {
	v reflect.Value
}

var (
	_ Reader     = (*Slice)(nil)
	_ SafeReader = (*Slice)(nil)
	_ ListerTo   = (*Slice)(nil)
)

func (s *Slice) Type() Type {
	return SliceType
}

func (s *Slice) Get(key string) (any, bool) {
	v, err := s.SafeGet(key)
	return v, err == nil
}

func (s *Slice) SafeGet(key string) (any, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: err,
		}
	}
	if n < 0 || n >= s.v.Len() {
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: errors.New("out of bounds error"),
		}
	}

	switch v := s.v.Index(n); {
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

func (s *Slice) List() []string {
	keys := make([]string, 0, s.v.Len())
	s.ListTo(&keys)
	return keys
}

func (s *Slice) ListTo(keys *[]string) {
	for i := 0; i < s.v.Len(); i++ {
		*keys = append(*keys, strconv.Itoa(i))
	}
}

func slice(v any) any {
	switch u := v.(type) {
	/*case Map:
	for k, v := range u {
		u[k] = slice(v)
	}

	return u
	*/
	case map[string]any:
		for k, v := range u {
			u[k] = slice(v)
		}

		return u
	case Interface:
		for _, k := range u.List() {
			v, _ := u.Get(k)
			_ = u.Set(k, slice(v))
		}

		return u
	case []any:
		m := make(map[string]any, len(u))

		for i, v := range u {
			m[fmt.Sprint(i)] = slice(v)
		}

		return m
	default:
		return v
	}
}

func unslice(v any) any {
	switch u := v.(type) {
	case []any:
		for i, v := range u {
			u[i] = unslice(v)
		}

		return u
	/*case Map:
	return unsliceMap(u)*/
	case map[string]any:
		return unsliceMap(u)
	case Interface:
		return unsliceInterface(u)
	default:
		return v
	}
}

func unsliceMap(m map[string]any) any {
	max := 0

	for k := range m {
		n, err := strconv.Atoi(k)
		if err != nil {
			max = -1
			break
		}

		if n > max {
			max = n
		}
	}

	if max == -1 {
		for k, v := range m {
			m[k] = unslice(v)
		}

		return m
	}

	w := make([]any, max+1)

	for k, v := range m {
		n, _ := strconv.Atoi(k)
		w[n] = unslice(v)
	}

	return w
}

func unsliceInterface(obj Interface) any {
	var (
		keys = obj.List()
		max  = 0
	)

	for _, k := range keys {
		n, err := strconv.Atoi(k)
		if err != nil {
			max = -1
			break
		}

		if n > max {
			max = n
		}
	}

	if max == -1 {
		for _, k := range keys {
			v, _ := obj.Get(k)
			_ = obj.Set(k, unslice(v))
		}

		return obj
	}

	w := make([]any, max+1)

	for _, k := range keys {
		v, _ := obj.Get(k)
		n, _ := strconv.Atoi(k)
		w[n] = unslice(v)
	}

	return w
}

func nclone[T any](v []T) []T {
	if len(v) == 0 {
		return nil
	}

	vCopy := make([]T, len(v))
	copy(vCopy, v)

	return vCopy
}
