package mergelist

import (
	"net/url"
	"strconv"
)

type Elem interface {
	~string | ~*url.URL
}

type List []struct {
	Next *List
	Key  Key
	URL  *url.URL
	Map  Map
}

var (
// _ Reader = (*List)(nil)
// _ Meta   = (*List)(nil)
)

func (l List) Type() Type {
	return TypeSlice
}

func (l List) Len() int {
	return len(l)
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
