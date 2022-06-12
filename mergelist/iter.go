package mergelist

import (
	"context"
	"net/url"

	"rafal.dev/objects/types"
)

type Iter interface {
	types.Iter
	Resolve() error
	URL() *url.URL
}

func (l *List) Range() Iter {
	return &iter{
		queue: newQueue(l),
	}
}

type elm struct {
	l    *List
	k    int
	leaf bool
}

func newQueue(l *List) []elm {
	return []elm{{l: l, k: 0}}
}

type iter struct {
	it    elm
	queue []elm
	done  bool
	err   error
}

var _ Iter = (*iter)(nil)

func (it *iter) Next(ctx context.Context) bool {
	if it.err != nil || len(it.queue) == 0 {
		it.done = true
		return false
	}

	var (
		n = len(it.queue) - 1
		k int
	)

	it.it, it.queue = it.queue[n], it.queue[:n]
	k = it.it.k

	if it.it.k++; it.it.k < it.it.l.Len() {
		it.queue = append(it.queue, it.it)
	}

	if l := (*it.it.l)[k].Next; l != nil {
		it.queue = append(it.queue, elm{l: l, k: 0})
	} else {
		it.it.leaf = true
	}

	return true
}

func (it *iter) Parent() Reader {
	return nil
	// return it.it.l
}

func (it *iter) Leaf() bool {
	return it.it.leaf
}

func (it *iter) Key() Key {
	return nil
	// return it.it.key
}

func (it *iter) Value() any {
	return it.it.l
}

func (it *iter) URL() *url.URL {
	return nil
}

func (it *iter) Resolve() error {
	return nil // todo
}

func (it *iter) Err() error {
	if !it.done {
		return ErrNotDone
	}
	return it.err
}
