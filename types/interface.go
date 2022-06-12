package types

type Reader interface {
	Get(key string) (value any, ok bool)
	List() []string
	Type() Type
}

type SafeReader interface {
	SafeGet(key string) (value any, err error)
}

type ListerTo interface {
	ListTo(*[]string)
}

type Writer interface {
	Del(key string) (ok bool)
	Set(key string, value any) (previous bool)
	Put(key string, hint Type) Writer
}

type SafeWriter interface {
	SafeDel(key string) error
	SafeSet(key string, value any) (previous bool, err error)
	SafePut(key string, hint Type) (Writer, error)
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
