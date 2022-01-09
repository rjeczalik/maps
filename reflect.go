package objects

import (
	"fmt"
	"reflect"
)

func typeOf(v any, iface bool) (t reflect.Type) {
	for t = reflect.TypeOf(v); t.Kind() == reflect.Ptr || (iface && t.Kind() == reflect.Interface); {
		t = t.Elem()
	}

	return t
}

func valueOf(v any, iface bool) (vv reflect.Value) {
	for vv = reflect.ValueOf(v); vv.Kind() == reflect.Ptr || (iface && vv.Kind() == reflect.Interface); {
		vv = vv.Elem()
	}

	return vv
}

func assertxtype(v interface{}) {
	if v == nil {
		return
	}

	switch typeOf(v, true).Kind() {
	case reflect.Array, reflect.Struct:
		panicf("unsupported type error, use objects.Build instead: %T", v)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panicf("unsupported type error: %T", v)
	}
}

func panicf(format string, v ...interface{}) {
	panic(fmt.Errorf(format, v...))
}
