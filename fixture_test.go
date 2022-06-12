package objects_test

import "rafal.dev/objects"

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

type Pairs []struct {
	K string
	V any
}

func (p Pairs) append(it objects.Iter) Pairs {
	kv := make(Pairs, 1)[0]
	kv.K, kv.V = it.Key().String(), it.Value()

	return append(p, kv)
}

func pairX() Pairs {
	return Pairs{
		{K: "A.B.C.D.0", V: 1},
		{K: "A.B.D.0", V: 2},
		{K: "A.B.D.1", V: 2},
		{K: "A.C.D.0", V: 3},
		{K: "A.C.D.1", V: 3},
		{K: "A.C.D.2", V: 3},
		{K: "A.D.0", V: 4},
		{K: "A.D.1", V: 4},
		{K: "A.D.2", V: 4},
		{K: "A.D.3", V: 4},
		{K: "B.C.D.0", V: 5},
		{K: "B.C.D.1", V: 5},
		{K: "B.C.D.2", V: 5},
		{K: "B.C.D.3", V: 5},
		{K: "B.C.D.4", V: 5},
		{K: "B.D.0", V: 6},
		{K: "B.D.1", V: 6},
		{K: "B.D.2", V: 6},
		{K: "B.D.3", V: 6},
		{K: "B.D.4", V: 6},
		{K: "B.D.5", V: 6},
		{K: "C.D.0", V: 7},
		{K: "C.D.1", V: 7},
		{K: "C.D.2", V: 7},
		{K: "C.D.3", V: 7},
		{K: "C.D.4", V: 7},
		{K: "C.D.5", V: 7},
		{K: "C.D.6", V: 7},
		{K: "D.0", V: 8},
		{K: "D.1", V: 8},
		{K: "D.2", V: 8},
		{K: "D.3", V: 8},
		{K: "D.4", V: 8},
		{K: "D.5", V: 8},
		{K: "D.6", V: 8},
		{K: "D.7", V: 8},
	}
}
