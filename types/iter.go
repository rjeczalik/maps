package types

type Iter interface {
	Next() bool
	Parent() Reader
	Leaf() bool
	Key() Key
	Value() any
	Err() error
}
