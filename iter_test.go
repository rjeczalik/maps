package objects_test

import (
	"testing"

	"rafal.dev/objects"
	"rafal.dev/objects/internal/misc"

	"github.com/google/go-cmp/cmp"
)

func TestIter(t *testing.T) {
	r := objects.Make(newX())

	cases := map[string]struct {
		it   objects.Iter
		want Pairs
	}{
		"iter": {
			it:   objects.Walk(r),
			want: pairX(),
		},
		"revIter": {
			it:   objects.Reverse(objects.Walk(r)),
			want: misc.Reverse(pairX()),
		},
	}

	for name, cas := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				it   = cas.it
				want = cas.want
				got  = make(Pairs, 0, len(want))
			)

			for it.Next() {
				if it.Leaf() {
					got = got.append(it)
				}
			}

			if err := it.Err(); err != nil {
				t.Fatalf("Err()=%+v", err)
			}

			if !cmp.Equal(got, want) {
				t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
			}
		})
	}
}
