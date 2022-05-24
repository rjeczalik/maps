package objects

type Iter interface {
	Next() bool
	Parent() Reader
	Leaf() bool
	Key() Key
	Value() any
	Err() error
}

func Walk(r Reader) Iter {
	return &iter{
		queue: []elm{{parent: r, left: r.List()}},
	}
}

func ReverseWalk(r Reader) Iter {
	return &revIter{
		queue: []elm{{parent: r, left: r.List()}},
	}
}

type elm struct {
	parent Reader
	key    []string
	left   []string
	v      any
	leaf   bool
}

type iter struct {
	it    elm
	queue []elm
	err   error
}

var _ Iter = (*iter)(nil)

func (it *iter) Next() bool {
	if len(it.queue) == 0 {
		return false
	}

	var (
		n = len(it.queue) - 1
		k string
	)

	it.it, it.queue = it.queue[n], it.queue[:n]
	k, it.it.left = it.it.left[0], it.it.left[1:]

	if len(it.it.left) != 0 {
		it.queue = append(it.queue, it.it)
	}

	it.it.key = clone(it.it.key, k)

	if it.it.v, it.err = get(it.it.parent, it.it.key...); it.err != nil {
		return false
	}

	if r, ok := it.it.v.(Reader); ok {
		it.queue = append(it.queue, elm{parent: r, key: it.it.key, left: r.List()})
	} else {
		it.it.leaf = true
	}

	return true
}

func (it *iter) Parent() Reader {
	return it.it.parent
}

func (it *iter) Leaf() bool {
	return it.it.leaf
}

func (it *iter) Key() Key {
	return it.it.key
}

func (it *iter) Value() any {
	return it.it.v
}

func (it *iter) Err() error {
	return it.err
}

type revIter struct {
	it    elm
	queue []elm
	rev   []elm
	err   error
}

var _ Iter = (*revIter)(nil)

func (rit *revIter) Next() bool {
	for len(rit.queue) != 0 {
		var (
			n    = len(rit.queue) - 1
			k    string
			it   elm
			leaf bool
		)

		it, rit.queue = rit.queue[n], rit.queue[:n]
		k, it.left = it.left[0], it.left[1:]

		if len(it.left) != 0 {
			rit.queue = append(rit.queue, it)
		}

		it.key = clone(it.key, k)

		if it.v, rit.err = get(it.parent, it.key...); rit.err != nil {
			return false
		}

		if r, ok := it.v.(Reader); ok {
			rit.queue = append(rit.queue, elm{parent: r, key: it.key, left: r.List()})
		} else {
			leaf = true
		}

		rit.rev = append(rit.rev, elm{parent: it.parent, key: it.key, v: it.v, leaf: leaf})
	}

	if len(rit.rev) == 0 {
		return false
	}

	var (
		n = len(rit.rev) - 1
	)

	rit.it, rit.rev = rit.rev[n], rit.rev[:n]

	return true
}

func (rit *revIter) Parent() Reader {
	return rit.it.parent
}

func (rit *revIter) Leaf() bool {
	return rit.it.leaf
}

func (rit *revIter) Key() Key {
	return rit.it.key
}

func (rit *revIter) Value() any {
	return rit.it.v
}

func (rit *revIter) Err() error {
	return rit.err
}
