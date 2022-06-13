package objects

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Slice struct {
	v reflect.Value
}

var (
	_ Reader = (*Slice)(nil)
	_ Meta   = (*Slice)(nil)
)

func (s *Slice) Type() Type {
	return TypeSlice
}

func (s *Slice) Get(ctx context.Context, key string) (any, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: err,
		}
	}
	if n < 0 || n >= s.v.Len() {
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Err: errors.New("out of bounds error"),
		}
	}

	switch v := s.v.Index(n); {
	case !v.CanInterface():
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Got: v,
			Err: fmt.Errorf("cannot access value: %s", v.Type()),
		}
	default:
		return tryMake(v.Interface()), nil
	}
}

func (s *Slice) List(ctx context.Context) ([]string, error) {
	keys := make([]string, 0, s.v.Len())

	for i := 0; i < s.v.Len(); i++ {
		keys = append(keys, strconv.Itoa(i))
	}

	return keys, nil
}
