package types

import "context"

type PrefixedReader struct {
	Key Key
	R   Reader
}

type PrefixedWriter struct {
	Key Key
	W   Writer
}

type Prefixed struct {
	PrefixedReader
	PrefixedWriter
}

var (
	_    Reader    = PrefixedReader{}
	_    Writer    = PrefixedWriter{}
	_    Interface = Prefixed{}
	_, _ Meta      = PrefixedReader{}, PrefixedWriter{}
)

func PrefixReader(r Reader, keys ...string) PrefixedReader {
	// todo: canonical
	return PrefixedReader{
		Key: keys,
		R:   r,
	}
}

func PrefixWriter(w Writer, keys ...string) PrefixedWriter {
	// todo: canonical
	return PrefixedWriter{
		Key: keys,
		W:   w,
	}
}

func Prefix(iface Interface, keys ...string) Prefixed {
	return Prefixed{
		PrefixedReader{
			Key: keys,
			R:   iface,
		},
		PrefixedWriter{
			Key: keys,
			W:   iface,
		},
	}
}

func (pr PrefixedReader) List(ctx context.Context) ([]string, error) {
	r, err := pr.base(ctx, "List")
	if err != nil {
		return nil, err
	}

	return r.List(ctx)
}

func (pr PrefixedReader) Type() Type {
	if m, ok := pr.R.(Meta); ok {
		return m.Type()
	}
	return nil
}

func (pr PrefixedReader) Get(ctx context.Context, key string) (any, error) {
	r, err := pr.base(ctx, "Get")
	if err != nil {
		return nil, err
	}

	v, err := r.Get(ctx, key)
	if err != nil {
		return nil, &Error{
			Op:  "Get",
			Key: append(pr.Key, key),
			Got: r,
			Err: ErrNotFound,
		}
	}

	return v, nil
}

func (pr PrefixedReader) base(ctx context.Context, op string) (Reader, error) {
	var (
		r  = pr.R
		ok bool
	)

	for i, key := range pr.Key {
		v, err := r.Get(ctx, key)
		if err != nil {
			return nil, &Error{
				Op:  op,
				Key: pr.Key[:i+1],
				Got: r,
				Err: err,
			}
		}

		if r, ok = v.(Reader); !ok {
			return nil, &Error{
				Op:   op,
				Key:  pr.Key[:i+1],
				Got:  v,
				Want: Reader(nil),
				Err:  ErrUnexpectedType,
			}
		}
	}

	return r, nil
}

func (pr PrefixedReader) reader() Reader {
	const maxDepth = 128

	key := pr.Key

	for i := 0; i < maxDepth; i++ {
		switch x := pr.R.(type) {
		case PrefixedReader:
			pr = x
			key.Prepend(pr.Key)
		case Prefixed:
			pr = x.PrefixedReader
			key.Prepend(pr.Key)
		default:
			return pr.R
		}
	}

	return nil
}

func (pw PrefixedWriter) Del(ctx context.Context, key string) error {
	pr, err := pw.reader("Del")
	if err != nil {
		return err
	}

	r, err := pr.base(ctx, "Del")
	if err != nil {
		return err
	}

	w, ok := r.(Writer)
	if !ok {
		return &Error{
			Op:   "Del",
			Key:  append(pw.Key, key),
			Got:  r,
			Want: Writer(nil),
			Err:  ErrUnexpectedType,
		}
	}

	if err := w.Del(ctx, key); err != nil {
		return &Error{
			Op:  "Del",
			Key: append(pw.Key, key),
			Err: err,
		}
	}

	return nil
}

func (pw PrefixedWriter) Set(ctx context.Context, key string, value any) error {
	pr, err := pw.reader("Set")
	if err != nil {
		return err
	}

	r, err := pr.base(ctx, "Set")
	if err != nil {
		return err
	}

	w, ok := r.(Writer)
	if !ok {
		return &Error{
			Op:   "Set",
			Key:  append(pw.Key, key),
			Got:  r,
			Want: Writer(nil),
			Err:  ErrUnexpectedType,
		}
	}

	if err := w.Set(ctx, key, value); err != nil {
		return &Error{
			Op:  "Set",
			Key: append(pw.Key, key),
			Err: err,
		}
	}

	return nil
}

func (pw PrefixedWriter) Put(ctx context.Context, key string, typ Type) (Writer, error) {
	var (
		w, k    = pw.writer()
		normkey = append(k, key)
		err     error
	)

	for i, key := range normkey {
		if w, err = w.Put(ctx, key, typ); err != nil {
			return nil, &Error{
				Op:  "Put",
				Key: normkey[:i+1],
				Got: w,
				Err: err,
			}
		}
	}

	return w, nil
}

func (pw PrefixedWriter) Type() Type {
	if m, ok := pw.W.(Meta); ok {
		return m.Type()
	}
	return nil
}

func (pw PrefixedWriter) reader(op string) (PrefixedReader, error) {
	w, key := pw.writer()
	r, ok := w.(Reader)
	if !ok {
		return PrefixedReader{}, &Error{
			Op:   op,
			Key:  pw.Key,
			Got:  pw.W,
			Want: Reader(nil),
			Err:  ErrUnexpectedType,
		}
	}

	return PrefixedReader{Key: key, R: r}, nil
}

func (pw PrefixedWriter) writer() (Writer, Key) {
	const maxDepth = 128

	key := pw.Key

	for i := 0; i < maxDepth; i++ {
		switch x := pw.W.(type) {
		case PrefixedWriter:
			pw = x
			key.Prepend(pw.Key)
		case Prefixed:
			pw = x.PrefixedWriter
			key.Prepend(pw.Key)
		default:
			return pw.W, key
		}
	}

	return nil, nil
}
