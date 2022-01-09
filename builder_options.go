package objects

import (
	"reflect"
	"strings"
)

func ConvSlice(v any) (any, error) {
	return v, nil
}

func FieldTag(tags ...string) FieldFunc {
	return func(f reflect.StructField) (name string, omit bool) {
		for _, tag := range tags {
			name, omit = f.Tag.Get(tag), false

			if i := strings.IndexRune(name, ','); i != -1 {
				for j, k := i+1, 0; j < len(name); j += k + 1 {
					if k = strings.IndexRune(name[j:], ','); k == -1 {
						k = len(name) - j
					}

					if name[j:j+k] == "omitempty" {
						omit = true
						break
					}
				}

				name = name[:i]
			}

			if name != "" && name != "-" {
				return name, omit
			}
		}

		return f.Name, false
	}
}

func MapSlice(v any) (any, error) {
	return nil, nil
}
