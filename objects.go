package objects

// TODO:
// - finish Del
// - finish Map
// - finish Slice
// - finish Struct
// - implement Zipper

func Copy(w Writer, r Reader) error {
	var (
		tr = TeeReader(r, w)
		it = Walk(tr)
	)

	for it.Next() {
	}

	return it.Err()
}

func get(r Reader, keys ...string) (v any, err error) {
	var n = len(keys) - 1

	if sr, ok := r.(SafeReader); ok {
		if v, err = sr.SafeGet(keys[n]); err != nil {
			return nil, &Error{
				Op:  "Get",
				Key: keys,
				Err: err,
			}
		}
	} else if v, ok = r.Get(keys[n]); !ok {
		return nil, &Error{
			Op:  "Get",
			Key: keys,
			Err: ErrNotFound,
		}
	}

	return v, nil
}

func set(w Writer, v any, keys ...string) (ok bool, err error) {
	var n = len(keys) - 1

	if sw, ok := w.(SafeWriter); ok {
		if ok, err = sw.SafeSet(keys[n], v); err != nil {
			return false, &Error{
				Op:  "Set",
				Key: keys,
				Err: err,
			}
		}

		return ok, nil
	}

	return w.Set(keys[n], v), nil
}

func del(w Writer, keys ...string) (err error) {
	var n = len(keys) - 1

	if sw, ok := w.(SafeWriter); ok {
		if err := sw.SafeDel(keys[n]); err != nil {
			return &Error{
				Op:  "Del",
				Key: keys,
				Err: err,
			}
		}

		return nil
	}

	if ok := w.Del(keys[n]); !ok {
		return &Error{
			Op:  "Del",
			Key: keys,
			Err: ErrNotFound,
		}
	}

	return nil
}

func put(w Writer, hint Type, keys ...string) (Writer, error) {
	var n = len(keys) - 1

	if sw, ok := w.(SafeWriter); ok {
		var err error
		if w, err = sw.SafePut(keys[n], hint); err != nil {
			return nil, &Error{
				Op:  "Put",
				Key: keys,
				Err: err,
			}
		}

		return w, nil
	}

	return w.Put(keys[n], hint), nil
}

func Get[T any](r Reader, keys ...string) (T, error) {
	var (
		t  T
		ok bool
	)

	for i := range keys[:len(keys)-1] {
		v, err := get(r, keys[:i+1]...)
		if err != nil {
			return t, err
		}

		if r, ok = v.(Reader); !ok {
			return t, &Error{
				Op:   "Get",
				Key:  keys[:i+1],
				Got:  v,
				Want: Reader(nil),
				Err:  ErrUnexpectedType,
			}
		}
	}

	v, err := get(r, keys...)
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

func Set(obj Interface, v any, keys ...string) (bool, error) {
	var n = len(keys) - 1

	w, err := Get[Writer](obj, keys[:n]...)
	if err != nil {
		return false, &Error{
			Op:  "Set",
			Key: keys[:n],
			Err: err,
		}
	}

	return set(w, v, keys...)
}

func Put(w Writer, hint Type, keys ...string) (Writer, error) {
	for i := range keys {
		if r, ok := w.(Reader); ok {
			if v, err := get(r, keys[:i+1]...); err == nil {
				if v, ok := v.(Writer); ok {
					w = v
					continue
				}

				return nil, &Error{
					Op:   "Put",
					Key:  keys[:i+1],
					Got:  v,
					Want: Writer(nil),
					Err:  ErrUnexpectedType,
				}
			}
		}

		var err error
		if w, err = put(w, hint, keys[:i+1]...); err != nil {
			return nil, err
		}
	}

	return w, nil
}

func Del(obj Interface, keys ...string) error {
	var n = len(keys) - 1

	w, err := Get[Writer](obj, keys[:n]...)
	if err != nil {
		return &Error{
			Op:  "Del",
			Key: keys,
			Err: err,
		}
	}

	return del(w, keys...)
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
