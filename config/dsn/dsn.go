package dsn

import (
	"net/url"
	"reflect"
	"strings"

	valid "github.com/go-playground/validator/v10"
)

var validator = valid.New()

// DSN Data Source Name
type DSN struct {
	*url.URL
}

// Parse parses a dsn string to DSN object
func Parse(dsn string) (*DSN, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	return &DSN{URL: u}, nil
}

// Bind binds DSN to specified value
func (d *DSN) Bind(v interface{}) error {
	vt := reflect.ValueOf(v)
	if vt.Kind() != reflect.Ptr || vt.IsNil() {
		return &BindTypeError{Type: reflect.TypeOf(v)}
	}
	if err := d.bind(vt); err != nil {
		return err
	}
	return validator.Struct(v)
}

func (d *DSN) bind(v reflect.Value) (err error) {
	afs := make(map[string]assignFunc)
	if d.User != nil {
		un := d.User.Username()
		afs["username"] = assignStringFunc(un)
		if pwd, ok := d.User.Password(); ok {
			afs["password"] = assignStringFunc(pwd)
		}
	}
	afs["address"] = assignAddressFunc(d.Addresses())
	afs["scheme"] = assignStringFunc(d.Scheme)
	dec := &decoder{
		query:       d.Query(),
		assignFuncs: afs,
	}
	return dec.decode(v)
}

func (d *DSN) Addresses() []string {
	switch d.Scheme {
	case "unix", "unixgram", "unixpacket":
		return []string{d.Path}
	}
	return strings.Split(d.Host, ",")
}
