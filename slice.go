package objects

import (
	"fmt"
	"strconv"
)

func slice(v any) any {
	switch u := v.(type) {
	case map[string]any:
		for k, v := range u {
			u[k] = slice(v)
		}

		return u
	case []any:
		m := make(map[string]any, len(u))

		for i, v := range u {
			m[fmt.Sprint(i)] = slice(v)
		}

		return m
	default:
		return v
	}
}

func unslice(v any) any {
	switch u := v.(type) {
	case []any:
		for i, v := range u {
			u[i] = unslice(v)
		}

		return u
	case map[string]any:
		max := 0

		for k := range u {
			n, err := strconv.Atoi(k)
			if err != nil {
				max = -1
				break
			}

			if n > max {
				max = n
			}
		}

		if max == -1 {
			for k, v := range u {
				u[k] = unslice(v)
			}

			return u
		}

		w := make([]any, max+1)

		for k, v := range u {
			n, _ := strconv.Atoi(k)
			w[n] = unslice(v)
		}

		return w
	default:
		return v
	}
}

func clone[T any](v []T) []T {
	if len(v) == 0 {
		return nil
	}

	vCopy := make([]T, len(v))
	copy(vCopy, v)

	return vCopy
}
