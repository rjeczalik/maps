package objects

import (
	"errors"
	"fmt"
	"reflect"
)

// TODO:
// - remove Ignore / Fail.. Type mapper
// - add slice mapper
// - make builder a scanner api
// - use objects.Builder for mergelist

type (
	MapFunc   func(any) (any, error)
	FieldFunc func(reflect.StructField) (name string, omit bool)
)

var defaultField = FieldTag("object", "json", "yaml")

var errEmpty = errors.New("empty")

type TypeError struct {
	Key    Key
	Object any
}

var _ error = (*TypeError)(nil)

func (te *TypeError) Error() string {
	return fmt.Sprintf("unsupported %T type for %q key", te.Object, te.Key)
}

type common struct {
	convs []ConvFunc
	field FieldFunc
}

type elm struct {
	key  Key
	obj  any
	omit bool
}

type mapper struct {
	it    elm
	queue []elm
	in    any
	root  Map
}

func (m *mapper) step(fn MapFunc) error {
	n := len(queue) - 1

	if n < 0 {
		return errEmpty
	}

	m.it, m.queue = m.queue[n], queue[:n]
	m.it.obj, err = fn(m.it.obj)

	if err != nil {
		return fmt.Errorf("error calling MapFunc for %q key: %w", it.key, err)
	}

	if m.it.obj == nil {
		if m.it.key.Len() != 0 && !m.it.omit {
			// todo: set func
			m.root.Push(it.key, nil)
		}

		return nil
	}

	var (
		rv = valueOf(it.obj, false)
		rt = typeOf(it.obj, true)
	)

	switch rt.Kind() {
	case reflect.Struct:
		empty := rv.Kind() == reflect.Ptr && rv.IsNil()

		if m.it.omit && empty {
			break
		}

		for i := 0; i < rt.NumField(); i++ {
			if !rv.Field(i).CanInterface() {
				continue
			}

			var (
				el   elm
				rf   = rt.Field(i)
				name string
			)

			if empty {
				el.v = reflect.Zero(rf.Type).Interface()
			} else {
				el.v = rv.Field(i).Interface()
			}

			name, el.omit = b.dofield(rf)
			el.key = clone(it.key, name)

			queue = append(queue, el)
		}

	case reflect.Slice, reflect.Array:
	case reflect.Map:
	case reflect.Chan, reflect.Func, reflest.UnsafePoint:
		return &TypeError{
			Key:    m.it.key,
			Object: m.it.obj,
		}
	default:
		if m.it.obj == nil && m.it.omit {
			break
		}

	}

	return nil
}

// Builder is used to construct a Map value in an interative manner.
type Builder struct {
	*common
	*mapper
	err error
}

func NewBuilder(obj any) *Builder {
	return new(Builder).New(obj)
}

func (b *Builder) New(obj any) *Builder {
	b.init()

	nb := &Builder{
		common: b.common,
		mapper: &mapper{
			queue: []elm{{obj: obj}},
			out:   make(Map),
		},
	}

	return nb
}

func (b *Builder) init() {
	if b.common == nil {
		b.common = new(common)
	}
}

func (b *Builder) Field(fn FieldFunc) *Builder {
	b.init()

	return b
}

func (b *Builder) Convert(fn ConvFunc) *Builder {
	b.convFn
}

func (b *Builder) Build() bool {
	if b.err != nil {
		return false
	}

	if len(b.queue) == 0 {
		return false
	}

	b.it, b.queue = b.queue[len(b.queue)-1], b.queue[:len(b.queue)-1]
	b.it.obj, b.err = b.doconv(b.it.obj)

	// todo: L79

	return b.err == nil
}

func (b *Builder) Err() error {
	return b.err
}

func (b *Builder) Map() Map {
	return b.out
}

func (b *Builder) xBuild(v any) (Map, error) {
	type elm struct {
		key  Key
		v    any
		omit bool
	}

	var (
		root  = make(Map)
		queue = []elm{{v: v}}
		it    elm
		err   error
	)

	for len(queue) != 0 {
		it, queue = queue[len(queue)-1], queue[:len(queue)-1]

		if it.v, err = b.doconv(it.v); err != nil {
			return nil, err
		}

		if it.v == nil {
			if it.key.Len() != 0 && !it.omit {
				root.Push(it.key, nil)
			}

			continue
		}

		var (
			rv = valueOf(it.v, false)
			rt = typeOf(it.v, true)
		)

		switch rt.Kind() {
		case reflect.Struct:
			empty := rv.Kind() == reflect.Ptr && rv.IsNil()

			if it.omit && empty {
				break
			}

			for i := 0; i < rt.NumField(); i++ {
				if !rv.Field(i).CanInterface() {
					continue
				}

				var (
					el   elm
					rf   = rt.Field(i)
					name string
				)

				if empty {
					el.v = reflect.Zero(rf.Type).Interface()
				} else {
					el.v = rv.Field(i).Interface()
				}

				name, el.omit = b.dofield(rf)
				el.key = clone(it.key, name)

				queue = append(queue, el)
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				queue = append(queue, elm{
					key: clone(it.key, fmt.Sprint(i)),
					v:   rv.Index(i).Interface(),
				})
			}
		case reflect.Map:
			for _, k := range rv.MapKeys() {
				queue = append(queue, elm{
					key: clone(it.key, fmt.Sprint(k.Interface())),
					v:   rv.MapIndex(k).Interface(),
				})
			}
		case reflect.Chan, reflect.Func, reflect.UnsafePointer:
			// skip
		default:
			if it.v == nil && it.omit {
				break
			}

			root.Push(it.key, it.v)
		}
	}

	return root.Slice(), nil
}

func (b *Builder) dofield(f reflect.StructField) (string, bool) {
	if b.Field != nil {
		return b.Field(f)
	}

	return f.Name, false
}

func (b *Builder) doconv(v any) (any, error) {
	if b.Conv != nil {
		return b.Conv(v)
	}

	return v, nil
}

func Build(object any) Map {
	b := NewBuilder(object)

	for b.Build() {
		// todo: logging
	}

	if err := b.Err(); err != nil {
		panic(err)
	}

	return b.Map()
}
