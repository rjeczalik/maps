package mergelist_test

import (
	"context"
	"fmt"
	"testing"

	"rafal.dev/objects/mergelist"
)

func TestMake(t *testing.T) {
	var (
		ctx = context.Background()
	)

	l, err := mergelist.Make(ctx, `file://testdata/list.yaml`)
	if err != nil {
		t.Fatalf("Make()=%+v", err)
	}

	var it mergelist.Iter

	for it = l.Range(); it.Next(ctx); {
		fmt.Printf("it.URL()=%q\n", it.URL())
		fmt.Printf("it.Key()=%q\n", it.Key())
		fmt.Printf("it.Value()=%#v\n", it.Value())
	}

	if err := it.Err(); err != nil {
		t.Fatalf("Err()=%+v", err)
	}
}
