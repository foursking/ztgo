package dsn

import (
	"encoding"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const queryPrefix = "query."

type decoder struct {
	query       url.Values
	assignFuncs map[string]assignFunc
}

func (d *decoder) decode(v reflect.Value) error {
	var tu encoding.TextUnmarshaler
	tu, v = d.indirect(v)
	if tu != nil {
		return tu.UnmarshalText([]byte(d.query.Encode()))
	}
	if v.Kind() != reflect.Struct {
		return &AssignTypeError{Value: d.query.Encode(), Type: v.Type()}
	}
	tv := v.Type()
	for i := 0; i < tv.NumField(); i++ {
		fv := v.Field(i)
		field := tv.Field(i)
		to := newTag(field.Tag.Get(tagName))
		if to.Name == "-" {
			continue
		}
		if af, ok := d.assignFuncs[to.Name]; ok {
			if err := af(fv, &tag{}); err != nil {
				return err
			}
			continue
		}
		if !strings.HasPrefix(to.Name, queryPrefix) {
			continue
		}
		to.Name = to.Name[len(queryPrefix):]
		if err := d.value(fv, "", to); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) value(v reflect.Value, prefix string, t *tag) (err error) {
	key := combinekey(prefix, t)
	var tu encoding.TextUnmarshaler
	tu, v = d.indirect(v)
	if tu != nil {
		if val, ok := d.query[key]; ok {
			return tu.UnmarshalText([]byte(val[0]))
		}
		if t.Default != "" {
			return tu.UnmarshalText([]byte(t.Default))
		}
		return
	}
	switch v.Kind() {
	case reflect.Bool:
		err = d.valueBool(v, prefix, t)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if v.Kind() == reflect.Int64 && v.Type() == reflect.TypeOf(time.Duration(1)) {
			err = d.valueTimeDuration(v, prefix, t)
			break
		}
		err = d.valueInt64(v, prefix, t)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		err = d.valueUint64(v, prefix, t)
	case reflect.Float32, reflect.Float64:
		err = d.valueFloat64(v, prefix, t)
	case reflect.String:
		err = d.valueString(v, prefix, t)
	case reflect.Slice:
		err = d.valueSlice(v, prefix, t)
	case reflect.Struct:
		err = d.valueStruct(v, t)
	case reflect.Ptr:
		if !d.hasKey(combinekey(prefix, t)) {
			break
		}
		if !v.CanSet() {
			break
		}
		nv := reflect.New(v.Type().Elem())
		v.Set(nv)
		err = d.value(nv, prefix, t)
	}
	return
}

func combinekey(prefix string, t *tag) string {
	key := t.Name
	if prefix != "" {
		key = prefix + "." + key
	}
	return key
}

func (d *decoder) hasKey(key string) bool {
	for k := range d.query {
		if strings.HasPrefix(k, key+".") || k == key {
			return true
		}
	}
	return false
}

func (d *decoder) valueBool(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setBool(v, val)
}

func (d *decoder) setBool(v reflect.Value, val string) error {
	bval, err := strconv.ParseBool(val)
	if err != nil {
		return &AssignTypeError{Value: val, Type: v.Type()}
	}
	v.SetBool(bval)
	return nil
}

func (d *decoder) valueInt64(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setInt64(v, val)
}

func (d *decoder) setInt64(v reflect.Value, val string) error {
	ival, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return &AssignTypeError{Value: val, Type: v.Type()}
	}
	v.SetInt(ival)
	return nil
}

func (d *decoder) valueTimeDuration(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setTimeDuration(v, val)
}

func (d *decoder) setTimeDuration(v reflect.Value, val string) error {
	dur, err := time.ParseDuration(val)
	if err != nil {
		return &AssignTypeError{Value: val, Type: v.Type()}
	}
	v.SetInt(int64(dur))
	return nil
}

func (d *decoder) valueUint64(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setUint64(v, val)
}

func (d *decoder) setUint64(v reflect.Value, val string) error {
	uival, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return &AssignTypeError{Value: val, Type: v.Type()}
	}
	v.SetUint(uival)
	return nil
}

func (d *decoder) valueFloat64(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setFloat64(v, val)
}

func (d *decoder) setFloat64(v reflect.Value, val string) error {
	fval, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return &AssignTypeError{Value: val, Type: v.Type()}
	}
	v.SetFloat(fval)
	return nil
}

func (d *decoder) valueString(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	val := d.query.Get(key)
	if val == "" {
		if t.Default == "" {
			return nil
		}
		val = t.Default
	}
	return d.setString(v, val)
}

func (d *decoder) setString(v reflect.Value, val string) error {
	v.SetString(val)
	return nil
}

func (d *decoder) valueSlice(v reflect.Value, prefix string, t *tag) error {
	key := combinekey(prefix, t)
	strs, ok := d.query[key]
	if !ok {
		strs = strings.Split(t.Default, ",")
	}
	if len(strs) == 0 {
		return nil
	}
	et := v.Type().Elem()
	var setFunc func(reflect.Value, string) error
	switch et.Kind() {
	case reflect.Bool:
		setFunc = d.setBool
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		setFunc = d.setInt64
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		setFunc = d.setUint64
	case reflect.Float32, reflect.Float64:
		setFunc = d.setFloat64
	case reflect.String:
		setFunc = d.setString
	default:
		return &AssignTypeError{Type: et, Value: strs[0]}
	}
	vals := reflect.MakeSlice(v.Type(), len(strs), len(strs))
	for i, str := range strs {
		if err := setFunc(vals.Index(i), str); err != nil {
			return err
		}
	}
	if v.CanSet() {
		v.Set(vals)
	}
	return nil
}

func (d *decoder) valueStruct(v reflect.Value, t *tag) error {
	tv := v.Type()
	for i := 0; i < tv.NumField(); i++ {
		fv := v.Field(i)
		field := tv.Field(i)
		fto := newTag(field.Tag.Get(tagName))
		if fto.Name == "-" {
			continue
		}
		if af, ok := d.assignFuncs[fto.Name]; ok {
			if err := af(fv, &tag{}); err != nil {
				return err
			}
			continue
		}
		if !strings.HasPrefix(fto.Name, queryPrefix) {
			continue
		}
		fto.Name = fto.Name[len(queryPrefix):]
		if err := d.value(fv, t.Name, fto); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) indirect(v reflect.Value) (encoding.TextUnmarshaler, reflect.Value) {
	v0 := v
	haveAddr := false

	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				haveAddr = false
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.Elem().Kind() != reflect.Ptr && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
				return u, reflect.Value{}
			}
		}
		if haveAddr {
			v = v0
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, v
}
