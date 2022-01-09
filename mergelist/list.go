package mergelist

import (
	"rafal.dev/objects"
)

type Func func(parent List, index int)

type List []struct {
	List List
	Map  objects.Map
}

func (l List) Set(key objects.Key, v any) error {
	return nil
}

func (l List) Merge() objects.Map {
	return nil // todo
}

func (l List) Walk(fn Func) {
	type elm struct {
		parent List
		keys   []int
	}

	var (
		queue = []elm{{parent: l, keys: l.keys()}}
		index int
		it    elm
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]
		index, it.keys = it.keys[0], it.keys[1:]

		fn(it.parent, index)

		if len(it.keys) != 0 {
			queue = append(queue, it)
		}

		if list := it.parent[index].List; len(list) != 0 {
			queue = append(queue, elm{parent: list, keys: list.keys()})
		}
	}
}

func (l List) keys() []int {
	n := make([]int, 0, len(l))

	for i := range l {
		n = append(n, i)
	}

	return n
}
