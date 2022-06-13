package types

import (
	"context"
	"sort"
)

type Map map[string]any

var (
	_ Interface = Map(nil)
	_ Meta      = Map(nil)
)

func (m Map) Type() Type {
	return TypeMap
}

func (m Map) Get(ctx context.Context, key string) (any, error) {
	v, ok := m[key]
	if !ok {
		return nil, &Error{
			Op:  "Get",
			Key: Key{key},
			Got: m,
			Err: ErrNotFound,
		}
	}
	return m.Type().Make(v), nil
}

func (m Map) List(ctx context.Context) ([]string, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (m Map) Del(ctx context.Context, key string) error {
	delete(m, key)
	return nil
}

func (m Map) Set(ctx context.Context, key string, value any) error {
	m[key] = value
	return nil
}

func (m Map) Put(ctx context.Context, key string, typ Type) (Writer, error) {
	switch x := m[key].(type) {
	case nil:
		if typ == nil {
			typ = m.Type()
		}

		w := typ.New()
		m[key] = w
		return w, nil
	case Writer:
		return x, nil
	default:
		return nil, &Error{
			Op:   "Put",
			Key:  Key{key},
			Got:  x,
			Want: Writer(nil),
			Err:  ErrUnexpectedType,
		}
	}
}
