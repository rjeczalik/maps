package types

import (
	"errors"
	"strconv"

	"golang.org/x/exp/slices"
)

type Slice []any

var (
	_ Interface  = (*Slice)(nil)
	_ ListerTo   = Slice(nil)
	_ SafeReader = Slice(nil)
	_ SafeWriter = (*Slice)(nil)
)

func (s Slice) Type() Type {
	return SliceType
}

func (s Slice) Get(key string) (any, bool) {
	v, err := s.SafeGet(key)
	return tryMake(v), err == nil
}

func (s Slice) List() []string {
	keys := make([]string, 0, len(s))
	s.ListTo(&keys)
	return keys
}

func (s Slice) ListTo(keys *[]string) {
	for i := range s {
		*keys = append(*keys, strconv.Itoa(i))
	}
}

func (s *Slice) Del(key string) bool {
	return s.SafeDel(key) == nil
}

func (s *Slice) Set(key string, value any) bool {
	previous, _ := s.SafeSet(key, value)
	return previous
}

func (s *Slice) Put(key string, hint Type) Writer {
	w, _ := s.SafePut(key, hint)
	return w
}

func (s Slice) SafeGet(key string) (value any, err error) {
	n, err := s.index(key, "Get")
	if err != nil {
		return nil, err
	}

	return tryMake(s[n]), nil
}

func (s *Slice) SafeDel(key string) error {
	n, err := s.index(key, "Del")
	if err != nil {
		return err
	}

	*s = slices.Delete(*s, n, n+1)

	return nil
}

func (s *Slice) SafeSet(key string, value any) (previous bool, err error) {
	n, err := s.index(key, "Set")
	if err != nil && (!errors.Is(err, ErrOutOfBounds) || n < 0) {
		return false, err
	}
	if m := len(*s); n >= m {
		*s = append(*s, make(Slice, n-m+1)...)
	} else {
		previous = true
	}

	(*s)[n] = value

	return previous, nil
}

func (s *Slice) SafePut(key string, hint Type) (Writer, error) {
	n, err := s.index(key, "Put")
	if err != nil && (!errors.Is(err, ErrOutOfBounds) || n < 0) {
		return nil, err
	}
	if m := len(*s); n >= m {
		*s = append(*s, make(Slice, n-m+1)...)
	} else if w, ok := tryMake((*s)[n]).(Writer); ok {
		return w, nil
	}

	var w Writer = makeOr(hint, &Slice{})
	(*s)[n] = w

	return w, nil
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
