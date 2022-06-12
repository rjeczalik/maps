package types

import (
	"context"
	"sort"
)

type Map map[string]any

var (
	_ Interface = Map(nil)
	_ ListerTo  = Map(nil)
)

func (m Map) Type() Type {
	return TypeMap
}

func (m Map) Get(ctx context.Context, key string) (any, bool) {
	v, ok := m[key]
	return tryMake(v), ok
}

func (m Map) List(ctx context.Context) []string {
	keys := make([]string, 0, len(m))
	m.ListTo(ctx, &keys)
	return keys
}

func (m Map) ListTo(ctx context.Context, keys *[]string) {
	for k := range m {
		*keys = append(*keys, k)
	}
	sort.Strings(*keys)
}

func (m Map) Del(ctx context.Context, key string) bool {
	_, ok := m[key]
	delete(m, key)
	return ok
}

func (m Map) Set(ctx context.Context, key string, value any) bool {
	_, ok := m[key]
	m[key] = value
	return ok
}

func (m Map) Put(ctx context.Context, key string, hint Type) Writer {
	v, ok := m[key].(Writer)
	if !ok {
		v = makeOr(hint, make(Map))
		m[key] = v
	}
	return v
}
