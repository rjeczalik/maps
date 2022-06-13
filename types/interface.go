package types

import "context"

var (
	TypeMap    Type = newMapType()
	TypeSlice  Type = newSliceType()
	TypeStruct Type = newStructType()
)

type Type interface {
	New() Interface
	Make(v any) any
	String() string
}

type Meta interface {
	Type() Type
}

func TypeOf(v any) Type {
	switch x := v.(type) {
	case Type:
		return x
	case Meta:
		return x.Type()
	default:
		return nil
	}
}

type Reader interface {
	Get(ctx context.Context, key string) (value any, err error)
	List(ctx context.Context) ([]string, error)
}

type Writer interface {
	Del(ctx context.Context, key string) error
	Set(ctx context.Context, key string, value any) error
	Put(ctx context.Context, key string, typ Type) (Writer, error)
}

type Interface interface {
	Reader
	Writer
}

type typ struct {
	name string
	fn   func() Interface
}

var _ Type = typ{}

func (t typ) String() string { return t.name }
func (t typ) New() Interface { return t.fn() }
func (t typ) Make(v any) any { return tryMake(v) }

func newMapType() Type {
	return typ{
		name: "Map",
		fn:   func() Interface { return Map{} },
	}
}

func newSliceType() Type {
	return typ{
		name: "Slice",
		fn:   func() Interface { return &Slice{} },
	}
}

func newStructType() Type {
	return typ{
		name: "Struct",
		fn:   func() Interface { return Map{} },
	}
}
