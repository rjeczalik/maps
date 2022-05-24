package types

import (
	"errors"
	"fmt"
)

var (
	ErrOutOfBounds    = errors.New("out of bounds")
	ErrNotFound       = errors.New("not found")
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
