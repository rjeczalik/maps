package codec

import "rafal.dev/objects"

type Codec interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

type codecFn struct {
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (fn codecFn) Marshal(v any) ([]byte, error) {
	return fn.marshal(v)
}

func (fn codecFn) Unmarshal(p []byte, v any) error {
	return fn.unmarshal(p, v)
}

type Map map[string]struct {
	Codec    Codec
	Priority int
	Children Map
}

func (m Map) Encode(key objects.Key, o objects.Object) ([]byte, error) {
	return nil, nil
}

func (m Map) Decode(key objects.Key, p []byte, o *objects.Object) error {
	return nil
}

func (m Map) Codec(key objects.Key) Codec {
	var (
		it = m[""]
		ok = false
	)

	it.Children = m

	for _, key := range key {
		if it, ok = it.Children[key]; !ok {
			return nil
		}
	}

	return it.Codec
}
