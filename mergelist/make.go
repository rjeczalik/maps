package mergelist

import (
	"fmt"
	"net/url"

	"rafal.dev/objects"
	"rafal.dev/objects/types"
)

func Make(v any) (List, error) {
	var (
		r = types.Make(v)
		l List
	)

	if r == nil {
		return l, l.append(v)
	}

	it := objects.Walk(r)

	for it.Next() {
		fmt.Println(it.Key(), it.Parent().Type())
	}

	if err := it.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (l *List) append(v any) error {
	s, ok := v.(string)
	if !ok {
		return &Error{
			Op:   "Make",
			Got:  v,
			Want: string(""),
			Err:  ErrUnexpectedType,
		}
	}
	if s == "" {
		return &Error{
			Op:  "Make",
			Err: ErrEmpty,
		}
	}

	u, err := url.Parse(s)
	if err != nil {
		return &Error{
			Op:  "Make",
			Got: s,
			Err: err,
		}
	}

	ll := make(List, 1)[0]
	ll.URL = u

	*l = append(*l, ll)

	return nil
}
