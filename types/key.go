package types

import "strings"

type Key []string

func (k Key) Base() string {
	if len(k) == 0 {
		return ""
	}

	return k[len(k)-1]
}

func (k Key) Dir() []string {
	if len(k) == 0 {
		return nil
	}

	return k[:len(k)-1]
}

func (k Key) Len() int {
	return len(k)
}

func (k Key) String() string {
	return strings.Join(k, ".")
}

func (k Key) Strings() []string {
	return k
}

func (k Key) Copy() Key {
	kCopy := make(Key, len(k))
	copy(kCopy, k)
	return kCopy
}

func (k *Key) Prepend(prefix Key) {
	n, m := len(*k), len(prefix)

	*k = append(*k, make([]string, m)...)

	for i := 0; i < n; i++ {
		(*k)[i+m] = (*k)[i]
	}

	copy(*k, prefix)
}

/*
type Pair struct {
	Key   Key
	Value any
}

type Pairs []Pair

func (p Pairs) Map() Map {
	m := make(Map)

	for _, p := range p {
		if err := Set(m, p.Value, p.Key...); err != nil {
			panic(err)
		}
	}

	return m
}*/
