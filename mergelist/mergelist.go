package mergelist

func nonempty(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}

	return ""
}

func clone(n []int, i ...int) []int {
	nCopy := make([]int, len(n)+len(i))
	copy(nCopy, n)
	copy(nCopy[len(n):], i)

	return nCopy
}
