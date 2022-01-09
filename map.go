package objects

import (
	"fmt"
	"sort"
	"strings"
)

type Func func(Map, Key)

type Map map[string]any

func Make[T interface{ map[string]any | []any }](v T) Map {
	var (
		root  = make(Map)
		queue = []elm{{v: v}}
		it    elm
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]

		assertxtype(it.v)

		switch v := it.v.(type) {
		case []any:
			for i, v := range v {
				queue = append(queue, elm{
					key: clone(it.key, fmt.Sprint(i)),
					v:   v,
				})
			}
		case map[string]any:
			for k, v := range v {
				queue = append(queue, elm{
					key: clone(it.key, k),
					v:   v,
				})
			}
		default:
			root.Push(it.key, it.v)
		}
	}

	return root.Slice()
}

func (m Map) Push(k Key, v any) (Map, Key) {
	if len(k) == 0 {
		return nil, nil
	}

	var parent = m

	for _, k := range k.Dir() {
		switch w := parent[k].(type) {
		case nil:
			child := make(map[string]any)
			parent[k] = child
			parent = child
		case map[string]any:
			parent = w
		case Map:
			parent = w
		default:
			return nil, nil
		}
	}

	switch parent[k.Base()].(type) {
	case map[string]any:
		return nil, nil
	case Map:
		return nil, nil
	}

	parent[k.Base()] = v

	return parent, k
}

func (m Map) Keys() []string {
	var keys = make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (m Map) longestKeys() []string {
	var keys = make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) >= len(keys[j])
	})

	return keys
}

// todo: rework as Child
func (m Map) Leaf(key Key) (any, bool) {
	v, ok := m[key.Base()]
	if ok {
		if _, ok = v.(map[string]any); ok {
			return nil, false
		}
	}

	return v, true
}

func (m Map) Child(keys ...string) (Map, bool) {
	var ok bool

	for _, k := range keys {
		if m, ok = m[k].(map[string]any); !ok {
			break
		}
	}

	return m, ok
}

func (m Map) Pairs() Pairs {
	var p Pairs

	m.Walk(func(m Map, k Key) {
		if v, ok := m.Leaf(k); ok {
			p = append(p, Pair{
				Key:   k,
				Value: v,
			})
		}
	})

	return p
}

func (m Map) Flat() Map {
	w := make(Map)

	m.Walk(func(m Map, k Key) {
		if v, ok := m.Leaf(k); ok {
			w[strings.Join(k, ".")] = v // todo: optimize allocs
		}
	})

	return w
}

func (m Map) Unflat() Map {
	w := make(Map)

	for k, v := range m {
		key := strings.Split(k, ".")

		w.Push(key, v)
	}

	return w
}

func (m Map) Slice() Map {
	// todo: remove the need for obj()
	return slice(m.obj()).(map[string]any)
}

func (m Map) Unslice() any {
	// todo: remove the need for obj()
	// todo: generics?
	return unslice(m.obj())
}

func (m Map) obj() map[string]any {
	return m
}

// todo: generics
func (m Map) Walk(fn Func) {
	if len(m) == 0 {
		return
	}

	type elm struct {
		parent Map
		key    []string
		left   []string
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: m, left: m.Keys()}}
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		key := clone(it.key, k) // todo: optimize allocs

		fn(it.parent, key)

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if child, ok := it.parent.Child(k); ok {
			queue = append(queue, elm{parent: child, key: key, left: child.Keys()})
		}
	}
}

func (m Map) ReverseWalk(fn Func) {
	if len(m) == 0 {
		return
	}

	type elm struct {
		parent Map
		key    []string
		left   []string
	}

	var (
		it    elm
		k     string
		queue = []elm{{parent: m, left: m.Keys()}}
		rev   []elm
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		k, it.left = it.left[0], it.left[1:]

		// todo: optimize allocs
		key := clone(it.key, k)

		rev = append(rev, elm{parent: it.parent, key: key})

		if len(it.left) != 0 {
			queue = append(queue, it)
		}

		if child, ok := it.parent.Child(k); ok {
			queue = append(queue, elm{parent: child, key: key, left: child.Keys()})
		}
	}

	for len(rev) != 0 {
		it, rev = rev[len(rev)-1], rev[:len(rev)-1]

		fn(it.parent, it.key)
	}
}
