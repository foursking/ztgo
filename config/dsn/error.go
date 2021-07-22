package dsn

import "reflect"

type BindTypeError struct {
	Type reflect.Type
}

func (e *BindTypeError) Error() string {
	if e.Type == nil {
		return "bind type is nil"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "bind type: " + e.Type.String() + " is not pointer"
	}
	return "invalid bind type: " + e.Type.String()
}

type AssignTypeError struct {
	Type  reflect.Type
	Value string
}

func (e *AssignTypeError) Error() string {
	return "cannot assign type " + e.Value + " to Go type " + e.Type.String()
}
