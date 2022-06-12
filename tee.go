package objects

import "context"

type teeReader struct {
	R Reader
	W Writer
}

var (
	_ Reader     = (*teeReader)(nil)
	_ SafeReader = (*teeReader)(nil)
	_ ListerTo   = (*teeReader)(nil)
)

func TeeReader(r Reader, w Writer) Reader {
	var (
		tr    = &teeReader{R: r, W: w}
		_, sr = r.(SafeReader)
		_, lt = r.(ListerTo)
	)

	switch {
	case sr && lt:
		return tr
	case sr:
		return struct {
			Reader
			SafeReader
		}{tr, tr}
	case lt:
		return struct {
			Reader
			ListerTo
		}{tr, tr}
	default:
		return struct{ Reader }{tr}
	}
}

func (tr *teeReader) Type() Type {
	return tr.R.Type()
}

func (tr *teeReader) Get(ctx context.Context, key string) (any, bool) {
	v, ok := tr.R.Get(ctx, key)
	if !ok {
		return nil, false
	}

	v, err := tr.tee(ctx, v, key)
	if err != nil {
		return nil, false
	}

	return v, true
}

func (tr *teeReader) SafeGet(ctx context.Context, key string) (any, error) {
	v, err := tr.R.(SafeReader).SafeGet(ctx, key)
	if err != nil {
		return nil, err
	}

	return tr.tee(ctx, v, key)
}

func (tr *teeReader) tee(ctx context.Context, v any, key string) (any, error) {
	if r, ok := tryMake(v).(Reader); ok {
		w, err := Put(ctx, tr.W, r.Type(), key)
		if err != nil {
			return nil, err
		}

		return TeeReader(r, w), nil
	}

	if _, err := Set(ctx, tr.W, v, key); err != nil {
		return nil, err
	}

	return v, nil
}

func (tr *teeReader) List(ctx context.Context) []string {
	return tr.R.List(ctx)
}

func (tr *teeReader) ListTo(ctx context.Context, keys *[]string) {
	tr.R.(ListerTo).ListTo(ctx, keys)
}
