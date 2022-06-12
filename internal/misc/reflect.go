package misc

import (
	"reflect"
)

func TypeOf(v any, iface bool) (t reflect.Type) {
	for t = reflect.TypeOf(v); t.Kind() == reflect.Ptr || (iface && t.Kind() == reflect.Interface); {
		t = t.Elem()
	}

	return t
}

func ValueOf(v any, iface bool) (vv reflect.Value) {
	for vv = reflect.ValueOf(v); vv.Kind() == reflect.Ptr || (iface && vv.Kind() == reflect.Interface); {
		vv = vv.Elem()
	}

	return vv
}
