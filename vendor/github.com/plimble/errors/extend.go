package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type _errorh struct {
	ErrMessage string `json:"message"`
	ErrCode    int    `json:"code"`
	*stack     `json:"-"`
}

func (e _errorh) Error() string { return e.ErrMessage }

func (e _errorh) Code() int { return e.ErrCode }

func (e _errorh) Message() string { return e.ErrMessage }

func (e _errorh) JSONError() error {
	b, _ := json.Marshal(e)
	return errors.New(string(b))
}

func (e _errorh) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.ErrMessage)
			fmt.Fprintf(s, "%+v", e.StackTrace())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.ErrMessage)
	}
}

// New returns an error with the supplied message.
func Newh(code int, message string) error {
	return _errorh{
		message,
		code,
		callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorhf(code int, format string, args ...interface{}) error {
	return _errorh{
		fmt.Sprintf(format, args...),
		code,
		callers(),
	}
}

func ParseJSON(err string) _errorh {
	e := _errorh{}
	errr := json.Unmarshal([]byte(err), &e)
	if errr != nil {
		e.ErrMessage = errr.Error()
		e.ErrCode = 500
		e.stack = callers()
	}

	e.stack = callers()
	return e
}

func BadRequest(message string) error {
	return &_errorh{
		message,
		400,
		callers(),
	}
}

func Unauthorized(message string) error {
	return &_errorh{
		message,
		401,
		callers(),
	}
}

func Forbidden(message string) error {
	return &_errorh{
		message,
		403,
		callers(),
	}
}

func NotFound(message string) error {
	return &_errorh{
		message,
		404,
		callers(),
	}
}

func InternalServerError(message string) error {
	return &_errorh{
		message,
		500,
		callers(),
	}
}

type httpStatus interface {
	Code() int
	Error() string
}

type microStatus interface {
	StatusCode() int32
	Error() string
}

func ErrorStatus(err error) (int, error) {
	if err == nil {
		return 0, nil
	}

	switch wrapErr := Cause(err).(type) {
	case httpStatus:
		return wrapErr.Code(), wrapErr
	case microStatus:
		return int(wrapErr.StatusCode()), wrapErr
	default:
		return 500, wrapErr
	}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	status, _ := ErrorStatus(err)
	if status == 404 {
		return true
	}

	return false
}

func IsInternalError(err error) bool {
	if err == nil {
		return false
	}
	status, _ := ErrorStatus(err)
	if status == 500 {
		return true
	}

	return false
}

func IsBadRequest(err error) bool {
	if err == nil {
		return false
	}
	status, _ := ErrorStatus(err)
	if status == 400 {
		return true
	}

	return false
}

func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	status, _ := ErrorStatus(err)
	if status == 401 {
		return true
	}

	return false
}

func IsForbidden(err error) bool {
	if err == nil {
		return false
	}
	status, _ := ErrorStatus(err)
	if status == 403 {
		return true
	}

	return false
}

func BadRequestErr(err error) error {
	if err == nil {
		return nil
	}

	return &_errorh{
		err.Error(),
		400,
		callers(),
	}
}

func UnauthorizedErr(err error) error {
	if err == nil {
		return nil
	}

	return &_errorh{
		err.Error(),
		401,
		callers(),
	}
}

func ForbiddenErr(err error) error {
	if err == nil {
		return nil
	}

	return &_errorh{
		err.Error(),
		403,
		callers(),
	}
}

func NotFoundErr(err error) error {
	if err == nil {
		return nil
	}

	return &_errorh{
		err.Error(),
		404,
		callers(),
	}
}

func InternalServerErrorErr(err error) error {
	if err == nil {
		return nil
	}

	return &_errorh{
		err.Error(),
		500,
		callers(),
	}
}
