package types_test

import (
	"testing"

	"rafal.dev/objects/types"

	"github.com/google/go-cmp/cmp"
)

func TestPrefixedReader(t *testing.T) {
	var (
		m  = newM()
		pr = types.PrefixReader(m, "foo", "bar")
	)

	if got, want := pr.Type(), types.TypeMap; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got := pr.List()
	want := []string{"dir", "file"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	v, ok := pr.Get("dir")
	if !ok {
		t.Fatalf("Get()=%t", ok)
	}

	r, ok := v.(types.Reader)
	if !ok {
		t.Fatalf("got %T, want %T", v, types.Reader(nil))
	}

	got = r.List()
	want = []string{"1", "2", "3"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	var (
		_, err = pr.SafeGet("notfound")
		e      = &types.Error{}
	)

	if !types.ErrAs(err, e, nil) {
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

	got = r.List()
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
		m  = newM()
		pw = types.PrefixWriter(m, "foo", "bar")
		pr = types.PrefixReader(m, "foo", "bar")
	)

	if err := pw.SafeDel("file"); err != nil {
		t.Fatalf("SafeDel()=%+v", err)
	}

	got := pr.List()
	want := []string{"dir"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	switch ok, err := pw.SafeSet("file", []byte("content")); {
	case err != nil:
		t.Fatalf("SafeSet()=%+v", err)
	case ok:
		t.Fatalf("got %t, want %t", ok, false)
	}

	got = pr.List()
	want = []string{"dir", "file"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	w, err := pw.SafePut("new", types.TypeMap)
	if err != nil {
		t.Fatalf("SafePut()=%+v", err)
	}

	_ = w.Set("a", 1)
	_ = w.Set("b", 2)
	_ = w.Set("c", 3)

	got = w.(types.Reader).List()
	want = []string{"a", "b", "c"}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s", cmp.Diff(got, want))
	}

	var (
		ppw = types.PrefixWriter(pw, "dir")
		ppr = types.PrefixReader(pr, "dir")
	)

	v, err := ppr.SafeGet("1")
	if err != nil {
		t.Fatalf("SafeGet()=%+v", err)
	}

	if v != 1 {
		t.Fatalf("got %#v, want %#v", v, 1)
	}

	switch ok, err := ppw.SafeSet("1", "foo"); {
	case err != nil:
		t.Fatalf("SafeSet()=%+v", err)
	case !ok:
		t.Fatalf("got %t, want %t", ok, true)
	}

	v, err = ppr.SafeGet("1")
	if err != nil {
		t.Fatalf("SafeGet()=%+v", err)
	}

	if v != "foo" {
		t.Fatalf("got %#v, want %#v", v, "foo")
	}
}
