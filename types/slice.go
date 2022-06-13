package types

import (
	"context"
	"errors"
	"strconv"

	"golang.org/x/exp/slices"
)

type Slice []any

var (
	_ Interface = (*Slice)(nil)
	_ Meta      = (*Slice)(nil)
)

func (s Slice) Type() Type {
	return TypeSlice
}

func (s Slice) Get(ctx context.Context, key string) (any, error) {
	n, err := s.index(key, "Get")
	if err != nil {
		return nil, err
	}

	return s.Type().Make(s[n]), nil
}

func (s Slice) List(ctx context.Context) ([]string, error) {
	keys := make([]string, 0, len(s))
	for i := range s {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys, nil
}

func (s *Slice) Del(ctx context.Context, key string) error {
	n, err := s.index(key, "Del")
	if err != nil {
		return err
	}

	*s = slices.Delete(*s, n, n+1)

	return nil
}

func (s *Slice) Set(ctx context.Context, key string, value any) error {
	n, err := s.index(key, "Set")
	if err != nil && (!errors.Is(err, ErrOutOfBounds) || n < 0) {
		return err
	}

	s.grow(n + 1)

	(*s)[n] = value

	return nil
}

func (s *Slice) Put(ctx context.Context, key string, typ Type) (Writer, error) {
	n, err := s.index(key, "Put")
	if err != nil && (!errors.Is(err, ErrOutOfBounds) || n < 0) {
		return nil, err
	}

	s.grow(n + 1)

	switch x := (*s)[n].(type) {
	case nil:
		if typ == nil {
			typ = s.Type()
		}

		w := typ.New()
		(*s)[n] = w
		return w, nil
	case Writer:
		return x, nil
	default:
		return nil, &Error{
			Op:   "Put",
			Key:  Key{key},
			Got:  x,
			Want: Writer(nil),
			Err:  ErrUnexpectedType,
		}
	}
}

func (s Slice) index(key, op string) (int, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return 0, &Error{
			Op:  op,
			Key: []string{key},
			Err: err,
		}
	}
	if n < 0 || n >= len(s) {
		return n, &Error{
			Op:   op,
			Key:  []string{key},
			Got:  n,
			Want: len(s),
			Err:  ErrOutOfBounds,
		}
	}

	return n, nil
}

func (s *Slice) grow(n int) {
	if m := len(*s); n > m {
		*s = append(*s, make(Slice, n-m)...)
	}
}
