package objects

type Field struct {
	Key   Key
	Value any
}

type Zipper struct {
}

type ZipFunc func(in Reader, kin string, out Interface, kout string) error

func Zip(in, out Interface) *Zipper {
	z := &Zipper{}

	return z
}

func (z *Zipper) Zip() bool {
	return true
}
