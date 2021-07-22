package errs

import (
	"encoding/json"
	"errors"
	"strconv"

	merr "github.com/micro/go-micro/v2/errors"
	pkgerr "github.com/pkg/errors"
)

//go:generate protoc -I. --go_out=paths=source_relative:. err.proto

// New creates error
func New(code int32, message string) *Err {
	return &Err{
		Code:    code,
		Message: message,
	}
}

// Error formats error as string
func (e *Err) Error() string {
	bs, _ := json.Marshal(e)
	return string(bs)
}

// Clone copies an error
func (e *Err) Clone() *Err {
	newErr := &Err{
		Code:    e.Code,
		Message: e.Message,
	}
	if e.Details != nil {
		newErr.Details = make(map[string]string)
		for k, v := range e.Details {
			newErr.Details[k] = v
		}
	}
	return newErr
}

// WithDetail adds error detail
// WithDetail("detail", "message")  Details is {"detail": "message"}
func (e *Err) WithDetail(v ...string) *Err {
	if e == nil {
		return nil
	}
	if len(v) == 0 {
		return e
	}
	newErr := e.Clone()
	if newErr.Details == nil {
		newErr.Details = make(map[string]string)
	}
	if len(v) == 1 {
		newErr.Details["detail"] = v[0]
		return newErr
	}
	newErr.Details[v[0]] = v[1]
	return newErr
}

// WithDetails appends error details
func (e *Err) WithDetails(m map[string]string) *Err {
	if e == nil {
		return nil
	}
	if m == nil {
		return e
	}
	newErr := e.Clone()
	if newErr.Details == nil {
		newErr.Details = make(map[string]string)
	}
	for k, v := range m {
		newErr.Details[k] = v
	}
	return newErr
}

// SetDetails sets error details
func (e *Err) SetDetails(m map[string]string) *Err {
	if e == nil {
		return nil
	}
	if m == nil {
		return e
	}
	newErr := e.Clone()
	newErr.Details = m
	return newErr
}

// Equal judges whether e1 equals to e2 by code
func Equal(e1, e2 *Err) bool {
	return e1.Code == e2.Code
}

// EqualError judges whether Err equals to error by code
func EqualError(e *Err, err error) bool {
	return e.Code == FromError(err).Code
}

// FromError unwraps Err from error, generates a new Err with error message if unwrap failed
func FromError(e error) *Err {
	if e == nil {
		return ErrOK
	}
	switch ec := pkgerr.Cause(e).(type) {
	case *Err:
		return ec
	case *merr.Error:
		return FromMicro(ec)
	}
	switch eu := errors.Unwrap(e).(type) {
	case *Err:
		return eu
	case *merr.Error:
		return FromMicro(eu)
	}
	return FromString(e.Error())
}

// FromString parse code string to error.
func FromString(errmsg string) *Err {
	if errmsg == "" {
		return ErrOK
	}
	i, err := strconv.Atoi(errmsg)
	if err == nil {
		return New(int32(i), errmsg)
	}
	return ErrInternalServerError.WithDetail("detail", errmsg)
}

func FromMicro(err *merr.Error) *Err {
	return ErrInternalServerError.WithDetails(map[string]string{
		"go_micro_code":   strconv.Itoa(int(err.Code)),
		"go_micro_status": err.Status,
		"go_micro_detail": err.Detail,
	})
}
