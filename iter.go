package objects

import "errors"

func Walk(r Reader) Iter {
	return &iter{
		queue: newQueue(r),
	}
}

func Reverse(it Iter) Iter {
	return &revIter{
		orig: it,
	}
}

type elm struct {
	parent Reader
	key    []string
	left   []string
	v      any
	leaf   bool
}

func newQueue(r Reader) []elm {
	return []elm{{parent: r, left: r.List()}}
}

type iter struct {
	it    elm
	queue []elm
	done  bool
	err   error
}

var _ Iter = (*iter)(nil)

func (it *iter) Next() bool {
	if it.err != nil || len(it.queue) == 0 {
		it.done = true
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

	if it.it.v, it.err = Get(it.it.parent, k); it.err != nil {
		it.done = true
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
	if !it.done {
		return ErrNotDone
	}
	return it.err
}

type revIter struct {
	orig Iter
	it   elm
	rev  []elm
	done bool
}

var _ Iter = (*revIter)(nil)

func (rit *revIter) Next() bool {
	for rit.orig.Next() {
		rit.rev = append(rit.rev, elm{
			parent: rit.orig.Parent(),
			leaf:   rit.orig.Leaf(),
			key:    rit.orig.Key(),
			v:      rit.orig.Value(),
		})
	}

	if rit.err() != nil {
		rit.done = true
		return false
	}

	if len(rit.rev) == 0 {
		rit.done = true
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
	if !rit.done {
		return ErrNotDone
	}
	return rit.err()
}

func (rit *revIter) err() (err error) {
	if err = rit.orig.Err(); errors.Is(err, ErrNotDone) {
		return nil
	}
	return err
}
