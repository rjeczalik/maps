package objects_test

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
