package objects

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

type Pair struct {
	Key   Key
	Value any
}

type Pairs []Pair

func (p Pairs) Map() Map {
	m := make(Map)

	for _, p := range p {
		if _, key := m.Push(p.Key, p.Value); key == nil {
			return nil
		}
	}

	return m
}
