package objects_test

import (
	"testing"

	"rafal.dev/objects"
	"rafal.dev/objects/simple"
)

func TestTeeReader(t *testing.T) {
	var (
		x  = newX()
		r  = objects.Make(x)
		w  = make(simple.Map)
		tr = objects.TeeReader(r, w)
		it = objects.Walk(tr)
	)

	_ = tr

	for it.Next() {
	}

	if err := it.Err(); err != nil {
		t.Fatalf("Err()=%+v", err)
	}

	if err := Equal(x, w); err != nil {
		t.Fatalf("Equal()=%s", err)
	}
}
