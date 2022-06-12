package objects_test

import (
	"context"
	"testing"

	"rafal.dev/objects"
	"rafal.dev/objects/types"
)

func TestTeeReader(t *testing.T) {
	var (
		x   = newX()
		r   = objects.Make(x)
		w   = make(types.Map)
		tr  = objects.TeeReader(r, w)
		it  = objects.Walk(tr)
		ctx = context.Background()
	)

	for it.Next(ctx) {
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
		x   = newX()
		r   = objects.Make(x)
		w   = make(types.Map)
		ctx = context.Background()
	)

	if err := objects.Copy(ctx, w, r); err != nil {
		t.Fatalf("Copy()=%+v", err)
	}

	if err := Equal(x, w); err != nil {
		t.Fatalf("Equal()=%s", err)
	}
}
