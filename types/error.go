package types

import (
	"errors"
	"fmt"
)

var (
	ErrOutOfBounds    = errors.New("index is out of bounds")
	ErrEmpty          = errors.New("unexpected empty value")
	ErrNotFound       = errors.New("not found")
	ErrNotDone        = errors.New("iterator not done")
	ErrUnexpectedType = errors.New("unexpected type")
)

type Error struct {
	Op   string
	Key  []string
	Got  any
	Want any
	Err  error
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	switch {
	case e.Got != nil && e.Want != nil:
		return fmt.Sprintf("%q operation error for %v key: got %#v, want %#v: %+v", e.Err, e.Key, e.Got, e.Want, e.Err)
	case e.Got != nil:
		return fmt.Sprintf("%q operation error for %v key and %#v value: %+v", e.Op, e.Key, e.Got, e.Err)
	case len(e.Key) != 0:
		return fmt.Sprintf("%q operation error for %v key: %+v", e.Op, e.Key, e.Err)
	default:
		return fmt.Sprintf("%q operation error: %+v", e.Op, e.Err)
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

func ErrAs(err error, out *Error, match func(*Error) bool) bool {
	const maxDepth = 128 // to prevent stack overflow if err has cyclic refs

	for i := 0; i < maxDepth; i++ {
		e := &Error{}

		switch ok := errors.As(err, &e); {
		case !ok:
			return false
		case ok && match == nil:
			*out = *e
			return true
		case ok && match(e):
			*out = *e
			return true
		default:
			err = e.Err
		}
	}

	fmt.Println("ErrAs not found")

	return false
}

func IsSentinelErr(err error) func(*Error) bool {
	return func(e *Error) bool {
		return e.Err == err
	}
}
