package objects_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "rafal.dev/objects"
)

type X struct {
	A struct {
		B struct {
			C struct {
				D []int
			}
			D []int
		}
		C struct {
			D []int
		}
		D []int
	}
	B struct {
		C struct {
			D []int
		}
		D []int
	}
	C struct {
		D []int
	}
	D []int
}

func newX() X {
	var x X

	x.A.B.C.D = []int{1}
	x.A.B.D = []int{2, 2}
	x.A.C.D = []int{3, 3, 3}
	x.A.D = []int{4, 4, 4, 4}
	x.B.C.D = []int{5, 5, 5, 5, 5}
	x.B.D = []int{6, 6, 6, 6, 6, 6}
	x.C.D = []int{7, 7, 7, 7, 7, 7, 7}
	x.D = []int{8, 8, 8, 8, 8, 8, 8, 8}

	return x
}

func TestBuilder(t *testing.T) {
	want, err := readmap("testdata/abcd.yaml")
	if err != nil {
		t.Fatal(err)
	}

	got := Build(newX())

	if !cmp.Equal(want, got) {
		t.Fatalf("want != got:\n%s\n", cmp.Diff(want, got))
	}
}

func makeField(tag string) reflect.StructField {
	return reflect.StructField{
		Tag: reflect.StructTag(tag),
	}
}

func TestFieldTag(t *testing.T) {
	tests := []struct {
		tags []string
		f    reflect.StructField
		name string
		omit bool
	}{
		0: {
			[]string{"json"},
			makeField(`json:"foo"`),
			"foo",
			false,
		},
		1: {
			[]string{"json", "yaml"},
			makeField(`yaml:"bar"`),
			"bar",
			false,
		},
		2: {
			[]string{"json", "yaml"},
			makeField(`json:"foo" yaml:"bar"`),
			"foo",
			false,
		},
		3: {
			[]string{"object", "json"},
			makeField(`json:"foo,omitempty"`),
			"foo",
			true,
		},
		4: {
			[]string{"object", "json"},
			makeField(`map:",omitempty" json:"foo,omitempty,string"`),
			"foo",
			true,
		},
		5: {
			[]string{"object", "json", "yaml"},
			makeField(`json:"-" map:",inline,omitempty" yaml:"foo,string,omitempty"`),
			"foo",
			true,
		}}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			name, omit := FieldTag(test.tags...)(test.f)

			if test.name != name {
				t.Errorf("got %q, want %q", name, test.name)
			}

			if test.omit != omit {
				t.Fatalf("got %t, want %t", omit, test.omit)
			}
		})
	}
}
