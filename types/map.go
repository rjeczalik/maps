package types

import "sort"

type Map map[string]any

var (
	_ Interface = Map(nil)
	_ ListerTo  = Map(nil)
)

func (m Map) Type() Type {
	return MapType
}

func (m Map) Get(key string) (any, bool) {
	v, ok := m[key]
	return tryMake(v), ok
}

func (m Map) List() []string {
	keys := make([]string, 0, len(m))
	m.ListTo(&keys)
	return keys
}

func (m Map) ListTo(keys *[]string) {
	for k := range m {
		*keys = append(*keys, k)
	}
	sort.Strings(*keys)
}

func (m Map) Del(key string) bool {
	_, ok := m[key]
	delete(m, key)
	return ok
}

func (m Map) Set(key string, value any) bool {
	_, ok := m[key]
	m[key] = value
	return ok
}

func (m Map) Put(key string, hint Type) Writer {
	v, ok := m[key].(Writer)
	if !ok {
		v = makeOr(hint, make(Map))
		m[key] = v
	}
	return v
}
