package objects_test

import (
	"fmt"
	"testing"

	"rafal.dev/objects"
)

func TestIter(t *testing.T) {
	var (
		x  = newX()
		r  = objects.Make(x)
		it = objects.Walk(r)
	)

	for it.Next() {
		if it.Leaf() {
			fmt.Printf("%v=%#v\n", it.Key(), it.Value())
		}
	}

	if err := it.Err(); err != nil {
		t.Fatalf("Err()=%+v", err)
	}
}

func TestRevIter(t *testing.T) {
	var (
		x   = newX()
		r   = objects.Make(x)
		rit = objects.ReverseWalk(r)
	)

	for rit.Next() {
		if rit.Leaf() {
			fmt.Printf("%v=%#v\n", rit.Key(), rit.Value())
		}
	}

	if err := rit.Err(); err != nil {
		t.Fatalf("Err()=%+v", err)
	}
}
