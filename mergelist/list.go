package mergelist

import (
	"errors"
	"net/url"
	"strconv"

	"golang.org/x/exp/slices"
)

type List []struct {
	Next *List
	Key  Key
	URL  *url.URL
	Map  Map
}

var (
	_ all = (*List)(nil)
)

func (l List) Type() Type {
	return TypeSlice
}

func (l List) Len() int {
	return len(l)
}

func (l List) Get(key string) (any, bool) {
	v, err := l.SafeGet(key)
	return v, err == nil
}

func (l List) SafeGet(key string) (any, error) {
	n, err := l.index(key, "Get")
	if err != nil {
		return nil, err
	}

	switch next, m, uri := l[n].Next, l[n].Map, l[n].URL; {
	case next != nil:
		return next, nil
	case m != nil:
		return m, nil
	case uri != nil:
		return uri, nil
	default:
		return nil, &Error{
			Op:  "Get",
			Key: []string{key},
			Got: l[n],
			Err: errors.New("empty list item"),
		}
	}
}

func (l List) List() []string {
	keys := make([]string, 0, l.Len())
	l.ListTo(&keys)
	return keys
}

func (l List) ListTo(keys *[]string) {
	for i := 0; i < l.Len(); i++ {
		*keys = append(*keys, strconv.Itoa(i))
	}
}

func (l *List) Del(key string) (ok bool) {
	return l.SafeDel(key) == nil
}

func (l *List) Set(key string, value any) (previous bool) {
	previous, _ = l.SafeSet(key, value)
	return previous
}

func (l *List) Put(key string, hint Type) Writer {
	w, _ := l.SafePut(key, hint)
	return w
}

func (l *List) SafeDel(key string) error {
	n, err := l.index(key, "Del")
	if err != nil {
		return err
	}

	*l = slices.Delete(*l, n, n+1)

	return nil
}

func (l *List) SafeSet(key string, value any) (previous bool, err error) {
	n, err := l.index(key, "Set")
	if err != nil && (!errors.Is(err, ErrOutOfBounds) || n < 0) {
		return false, err
	}

	if n >= l.Len() {
		*l = append(*l, make(List, n-l.Len()+1)...)
	} else {
		previous = true
	}

	switch v := value.(type) {
	case string:
		u, err := url.Parse(v)
		if err != nil {
			return false, &Error{
				Op:  "Set",
				Key: []string{key},
				Got: v,
				Err: err,
			}
		}

		(*l)[n].URL = u
	case []any:
		// todo
	case map[string]any:
		// todo
	}

	return false, nil
}

func (l *List) SafePut(key string, hint Type) (Writer, error) {
	n, err := l.index(key, "Put")
	if err != nil {
		return nil, err
	}

	_ = n // todo

	return nil, nil
}

func (l List) index(key, op string) (int, error) {
	n, err := strconv.Atoi(key)
	if err != nil {
		return 0, &Error{
			Op:  op,
			Key: []string{key},
			Err: err,
		}
	}
	if n < 0 || n >= l.Len() {
		return n, &Error{
			Op:   op,
			Key:  []string{key},
			Got:  n,
			Want: l.Len(),
			Err:  ErrOutOfBounds,
		}
	}

	return n, nil
}
