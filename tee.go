package objects

import (
	"context"

	"rafal.dev/objects/types"
)

type teeReader struct {
	R Reader
	W Writer
}

var (
	_ Reader = (*teeReader)(nil)
	_ Meta   = (*teeReader)(nil)
)

func TeeReader(r Reader, w Writer) Reader {
	return &teeReader{R: r, W: w}
}

func (tr *teeReader) Type() Type {
	if m, ok := tr.R.(Meta); ok {
		return m.Type()
	}
	return nil
}

func (tr *teeReader) Get(ctx context.Context, key string) (any, error) {
	v, err := tr.R.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	v, err = tr.tee(ctx, v, key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (tr *teeReader) tee(ctx context.Context, v any, key string) (any, error) {
	if r, ok := tryMake(v).(Reader); ok {
		w, err := tr.W.Put(ctx, key, types.TypeOf(r))
		if err != nil {
			return nil, err
		}

		return TeeReader(r, w), nil
	}

	if err := tr.W.Set(ctx, key, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (tr *teeReader) List(ctx context.Context) ([]string, error) {
	return tr.R.List(ctx)
}
