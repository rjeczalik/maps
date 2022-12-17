package mergelist

import (
	"context"

	"rafal.dev/objects/types"
)

type Iter interface {
	types.Iter
}

func (b *Builder) Range(it types.Iter) Iter {
	return nil
}

type iter struct {
}

var _ Iter = (*iter)(nil)

func (it *iter) Next(ctx context.Context) bool {
	return true
}

func (it *iter) Parent() Reader {
	return nil
}

func (it *iter) Leaf() bool {
	return true
}

func (it *iter) Key() Key {
	return nil
}

func (it *iter) Value() any {
	return nil
}

func (it *iter) Err() error {
	_ = ErrNotDone
	return nil
}
