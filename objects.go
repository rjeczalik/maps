package objects

import (
	"context"
	"errors"
)

func Copy(ctx context.Context, w Writer, r Reader) error {
	var (
		tr = TeeReader(r, w)
		it = Walk(tr)
	)

	for it.Next(ctx) {
	}

	return it.Err()
}

func TGet[T any](ctx context.Context, r Reader, keys ...string) (T, error) {
	var (
		t  T
		ok bool
	)

	v, err := Get(ctx, r, keys...)
	if err != nil {
		return t, err
	}

	if t, ok = v.(T); !ok {
		return t, &Error{
			Op:   "Get",
			Key:  keys,
			Got:  v,
			Want: t,
			Err:  ErrUnexpectedType,
		}
	}

	return t, nil
}

func Get(ctx context.Context, r Reader, keys ...string) (any, error) {
	var n = len(keys) - 1

	if n < 0 {
		return nil, &Error{
			Op:  "Get",
			Err: errors.New("keys are empty"),
		}
	}

	return PrefixedReader{
		Key: keys[:n],
		R:   r,
	}.Get(ctx, keys[n])
}

func Set(ctx context.Context, w Writer, v any, keys ...string) error {
	var n = len(keys) - 1

	if n < 0 {
		return &Error{
			Op:  "Set",
			Err: errors.New("keys are empty"),
		}
	}

	return PrefixedWriter{
		Key: keys[:n],
		W:   w,
	}.Set(ctx, keys[n], v)
}

func Put(ctx context.Context, w Writer, hint Type, keys ...string) (Writer, error) {
	var n = len(keys) - 1

	if n < 0 {
		return nil, &Error{
			Op:  "Put",
			Err: errors.New("keys are empty"),
		}
	}

	return PrefixedWriter{
		Key: keys[:n],
		W:   w,
	}.Put(ctx, keys[n], hint)
}

func Del(ctx context.Context, w Writer, keys ...string) error {
	var n = len(keys) - 1

	if n < 0 {
		return &Error{
			Op:  "Del",
			Err: errors.New("keys are empty"),
		}
	}

	return PrefixedWriter{
		Key: keys[:n],
		W:   w,
	}.Del(ctx, keys[n])

}

func clone(s []string, vs ...string) []string {
	sCopy := make([]string, len(s), len(s)+len(vs))
	copy(sCopy, s)

	for _, v := range vs {
		if v != "" && v != "-" {
			sCopy = append(sCopy, v)
		}
	}

	return sCopy
}

func npmin(i, j int) int {
	if i < j || j < 0 {
		return i
	}
	return j
}
