package types

import "context"

type Reader interface {
	Get(ctx context.Context, key string) (value any, ok bool)
	List(ctx context.Context) []string
	Type() Type
}

type SafeReader interface {
	SafeGet(ctx context.Context, key string) (value any, err error)
}

type ListerTo interface {
	ListTo(context.Context, *[]string)
}

type Writer interface {
	Del(ctx context.Context, key string) (ok bool)
	Set(ctx context.Context, key string, value any) (previous bool)
	Put(ctx context.Context, key string, hint Type) Writer
}

type SafeWriter interface {
	SafeDel(ctx context.Context, key string) error
	SafeSet(ctx context.Context, key string, value any) (previous bool, err error)
	SafePut(ctx context.Context, key string, hint Type) (Writer, error)
}

type Interface interface {
	Reader
	Writer
}

type SafeInterface interface {
	Interface
	SafeReader
	SafeWriter
}
