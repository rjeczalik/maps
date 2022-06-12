package types_test

import (
	"strconv"
	"testing"

	"rafal.dev/objects/types"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestErrFrom(t *testing.T) {
	opts := []cmp.Option{cmpopts.EquateErrors()}

	cases := []struct {
		err   error
		match func(*types.Error) bool
		want  *types.Error
	}{
		0: {
			err:   wrap(10, sentinel(types.ErrNotFound), got(9, "magic")),
			match: types.IsSentinelErr(types.ErrNotFound),
			want:  &types.Error{Op: "9", Got: "magic", Err: types.ErrNotFound},
		},
	}

	for _, cas := range cases {
		t.Run("", func(t *testing.T) {
			got := &types.Error{}

			switch ok := types.ErrAs(cas.err, got, cas.match); {
			case !ok && cas.want != nil:
				t.Errorf("error not found: %#v", cas.want)
			case ok && cas.want == nil:
				t.Errorf("unexpected error found: %#v", got)
			case ok && !cmp.Equal(*got, *cas.want, opts...):
				t.Errorf("got != want:\n%s", cmp.Diff(*got, *cas.want, opts...))
			}
		})
	}
}

func wrap(n int, fns ...func(int, *types.Error)) *types.Error {
	var (
		root = &types.Error{Op: "0"}
		it   = root
	)

	for i := 0; i < n-1; i++ {
		err := &types.Error{Op: strconv.Itoa(i + 1)}

		it.Err = err
		it = err
	}

	it = root

	for i := 0; i < n; i++ {
		for _, fn := range fns {
			fn(i, it)
		}

		if e, ok := it.Err.(*types.Error); ok {
			it = e
		} else {
			break
		}
	}

	return root
}

func sentinel(err error) func(int, *types.Error) {
	return func(i int, e *types.Error) {
		if e.Err == nil {
			e.Err = err
			e.Got = got
		}
	}
}

func got(i int, got any) func(int, *types.Error) {
	return func(j int, e *types.Error) {
		if j == i {
			e.Got = got
		}
	}
}
