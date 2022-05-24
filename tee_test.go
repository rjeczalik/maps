package objects_test

import (
	"testing"

	"rafal.dev/objects"
	"rafal.dev/objects/types"
)

func TestTeeReader(t *testing.T) {
	var (
		x  = newX()
		r  = objects.Make(x)
		w  = make(types.Map)
		tr = objects.TeeReader(r, w)
		it = objects.Walk(tr)
	)

	for it.Next() {
	}

	if err := it.Err(); err != nil {
		t.Fatalf("Err()=%+v", err)
	}

	if err := Equal(x, w); err != nil {
		t.Fatalf("Equal()=%s", err)
	}
}

func TestCopy(t *testing.T) {
	var (
		x = newX()
		r = objects.Make(x)
		w = make(types.Map)
	)

	if err := objects.Copy(w, r); err != nil {
		t.Fatalf("Copy()=%+v", err)
	}

	if err := Equal(x, w); err != nil {
		t.Fatalf("Equal()=%s", err)
	}
}
