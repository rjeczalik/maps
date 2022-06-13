package types_test

import (
	"context"
	"testing"

	"rafal.dev/objects/types"

	"github.com/google/go-cmp/cmp"
)

func TestPrefixedReader(t *testing.T) {
	var (
		m   = newM()
		pr  = types.PrefixReader(m, "foo", "bar")
		ctx = context.Background()
	)

	if got, want := pr.Type().String(), types.TypeMap.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got, _ := pr.List(ctx)
	want := []string{"dir", "file"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	v, err := pr.Get(ctx, "dir")
	if err != nil {
		t.Fatalf("Get()=%+v", err)
	}

	r, ok := v.(types.Reader)
	if !ok {
		t.Fatalf("got %T, want %T", v, types.Reader(nil))
	}

	got, _ = r.List(ctx)
	want = []string{"1", "2", "3"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	var e types.Error

	_, err = pr.Get(ctx, "notfound")
	if !types.ErrAs(err, &e, nil) {
		t.Fatalf("got %T, want %T", err, e)
	}

	if e.Err != types.ErrNotFound {
		t.Fatalf("got %#v, want %#v", e.Err, types.ErrNotFound)
	}

	got = e.Key
	want = []string{"foo", "bar", "notfound"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	r, ok = e.Got.(types.Reader)
	if !ok {
		t.Fatalf("got %T, want %T", e.Got, types.Reader(nil))
	}

	got, _ = r.List(ctx)
	want = []string{"dir", "file"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	_, ok = r.(types.Writer)
	if !ok {
		t.Fatalf("got %T, want %T", r, types.Writer(nil))
	}
}

func TestPrefixedWriter(t *testing.T) {
	var (
		m   = newM()
		pw  = types.PrefixWriter(m, "foo", "bar")
		pr  = types.PrefixReader(m, "foo", "bar")
		ctx = context.Background()
	)

	if err := pw.Del(ctx, "file"); err != nil {
		t.Fatalf("SafeDel()=%+v", err)
	}

	got, _ := pr.List(ctx)
	want := []string{"dir"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	if err := pw.Set(ctx, "file", []byte("content")); err != nil {
		t.Fatalf("SafeSet()=%+v", err)
	}

	got, _ = pr.List(ctx)
	want = []string{"dir", "file"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	w, err := pw.Put(ctx, "new", types.TypeMap)
	if err != nil {
		t.Fatalf("SafePut()=%+v", err)
	}

	_ = w.Set(ctx, "a", 1)
	_ = w.Set(ctx, "b", 2)
	_ = w.Set(ctx, "c", 3)

	got, _ = w.(types.Reader).List(ctx)
	want = []string{"a", "b", "c"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	var (
		ppw = types.PrefixWriter(pw, "dir")
		ppr = types.PrefixReader(pr, "dir")
	)

	v, err := ppr.Get(ctx, "1")
	if err != nil {
		t.Fatalf("SafeGet()=%+v", err)
	}

	if v != 1 {
		t.Fatalf("got %#v, want %#v", v, 1)
	}

	if err := ppw.Set(ctx, "1", "foo"); err != nil {
		t.Fatalf("SafeSet()=%+v", err)
	}

	v, err = ppr.Get(ctx, "1")
	if err != nil {
		t.Fatalf("SafeGet()=%+v", err)
	}

	if v != "foo" {
		t.Fatalf("got %#v, want %#v", v, "foo")
	}
}
