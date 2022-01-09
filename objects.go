package objects

func clone(s []string, vs ...string) []string {
	sCopy := make([]string, len(s), len(s)+len(vs))
	copy(sCopy, s)

	for _, v := range vs {
		if v != "" && v != "-" {
			sCopy = append(sCopy, v)
		}
	}

	return sCopy
}

func npmin(i, j int) int {
	if i < j || j < 0 {
		return i
	}
	return j
}
