package objects_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"rafal.dev/objects"

	"github.com/google/go-cmp/cmp"
	yaml "gopkg.in/yaml.v3"
)

func TestMake(t *testing.T) {
	m := make([]any, 2)
	_ = m

	_ = objects.Make(m)
}

func TestMapWalk(t *testing.T) {
	m, err := readmap("testdata/abcd.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var keys, rkeys []string

	m.Walk(func(m objects.Map, k objects.Key) {
		keys = append(keys, k.String())
	})

	m.ReverseWalk(func(m objects.Map, k objects.Key) {
		rkeys = append(rkeys, k.String())
	})

	if len(keys) != len(rkeys) {
		t.Fatalf("len(keys)=%d != len(rkeys)=%d", len(keys), len(rkeys))
	}

	for i, j := 0, len(keys)-1; i < len(keys); i, j = i+1, j-1 {
		if keys[i] != rkeys[j] {
			t.Fatalf("keys[%d]=%q != rkeys[%d]=%q", i, keys[i], j, rkeys[j])
		}
	}
}

func TestMap(t *testing.T) {
	want, err := readmap("testdata/abcd.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if got := want.Pairs().Map(); !cmp.Equal(want, got) {
		t.Fatalf("want != got:\n%s\n", cmp.Diff(want, got))
	}

	wflat := want.Flat()

	if *updateGolden {
		p, err := yaml.Marshal(wflat)
		if err != nil {
			t.Fatalf("yaml.Marshal()=%s", err)
		}

		if err := ioutil.WriteFile("testdata/abcd.flat.yaml.golden", p, 0644); err != nil {
			t.Fatalf("ioutil.WriteFile()=%s", err)
		}

		return
	}

	gflat, err := readmap("testdata/abcd.flat.yaml.golden")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wflat, gflat) {
		t.Fatalf("want != got:\n%s\n", cmp.Diff(wflat, gflat))
	}

	if wflat, gflat = wflat.Unslice().(map[string]any), gflat.Unslice().(map[string]any); !cmp.Equal(wflat, gflat) {
		t.Fatalf("want != got:\n%s\n", cmp.Diff(wflat, gflat))
	}
}

func BenchmarkMap(b *testing.B) {
	m, err := readmap("testdata/ansible-facts.json")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m
	}
}

func readmap(file string) (objects.Map, error) {
	p, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var v map[string]any

	if err := yaml.Unmarshal(p, &v); err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	return objects.Map(v).Slice(), nil
}
