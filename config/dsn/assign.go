package dsn

import (
	"reflect"
	"strings"
)

type assignFunc func(reflect.Value, *tag) error

func assignStringFunc(s string) assignFunc {
	return func(v reflect.Value, t *tag) error {
		if v.Kind() != reflect.String || !v.CanSet() {
			return &AssignTypeError{Value: "string", Type: v.Type()}
		}
		if s == "" {
			v.SetString(t.Default)
		} else {
			v.SetString(s)
		}
		return nil
	}
}

func assignAddressFunc(addrs []string) assignFunc {
	return func(v reflect.Value, t *tag) error {
		if v.Kind() == reflect.String {
			if addrs[0] == "" && t.Default != "" {
				v.SetString(t.Default)
			} else {
				v.SetString(addrs[0])
			}
			return nil
		}
		if !(v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String) {
			return &AssignTypeError{Value: strings.Join(addrs, ","), Type: v.Type()}
		}
		vals := reflect.MakeSlice(v.Type(), len(addrs), len(addrs))
		for i, addr := range addrs {
			vals.Index(i).SetString(addr)
		}
		if v.CanSet() {
			v.Set(vals)
		}
		return nil
	}
}
