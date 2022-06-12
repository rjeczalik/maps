package types_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"rafal.dev/objects/types"
)

func TestKeyPrepend(t *testing.T) {
	cases := []struct {
		orig   types.Key
		prefix types.Key
		want   types.Key
	}{
		0: {
			orig:   types.Key{},
			prefix: types.Key{"foo"},
			want:   types.Key{"foo"},
		},
		1: {
			orig:   types.Key{"foo", "bar"},
			prefix: types.Key{"baz", "qux"},
			want:   types.Key{"baz", "qux", "foo", "bar"},
		},
	}

	for _, cas := range cases {
		t.Run("", func(t *testing.T) {
			got := cas.orig.Copy()

			got.Prepend(cas.prefix)

			if !cmp.Equal(got, cas.want) {
				t.Errorf("got != want:\n%s", cmp.Diff(got, cas.want))
			}
		})
	}
}
