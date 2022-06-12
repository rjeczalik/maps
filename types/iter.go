package types

import "context"

type Iter interface {
	Next(context.Context) bool
	Parent() Reader
	Leaf() bool
	Key() Key
	Value() any
	Err() error
}
