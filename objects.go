package objects

import "errors"

func Copy(w Writer, r Reader) error {
	var (
		tr = TeeReader(r, w)
		it = Walk(tr)
	)

	for it.Next() {
	}

	return it.Err()
}

func TGet[T any](r Reader, keys ...string) (T, error) {
	var (
		t  T
		ok bool
	)

	v, err := Get(r, keys...)
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

func Get(r Reader, keys ...string) (any, error) {
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
	}.SafeGet(keys[n])
}

func Set(w Writer, v any, keys ...string) (bool, error) {
	var n = len(keys) - 1

	if n < 0 {
		return false, &Error{
			Op:  "Set",
			Err: errors.New("keys are empty"),
		}
	}

	return PrefixedWriter{
		Key: keys[:n],
		W:   w,
	}.SafeSet(keys[n], v)
}

func Put(w Writer, hint Type, keys ...string) (Writer, error) {
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
	}.SafePut(keys[n], hint)
}

func Del(w Writer, keys ...string) error {
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
	}.SafeDel(keys[n])

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
