package types_test

import "rafal.dev/objects/types"

type (
	M = types.Map
)

func newM() M {
	return M{
		"foo": M{
			"bar": M{
				"dir": M{
					"1": 1,
					"2": 2,
					"3": 3,
				},
				"file": []byte("content"),
			},
		},
	}
}
