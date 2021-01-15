package errors

import (
	"fmt"
	"net/http"
)

//Error Error
type Error struct {
	ECode   int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

//Error Error
func (e *Error) Error() string {
	return e.Message
}

//Code Code
func (e *Error) Code() int {
	return e.ECode
}

//HTTPError Error
type HTTPError struct {
	error
	Code int
}

//Status Status
func (e *HTTPError) Status() int {
	return e.Code
}

//New New
func New(msg string) error {
	return &Error{Message: msg}
}

//Code NewCode
func Code(code int, msg string) error {
	return &Error{ECode: code, Message: msg}
}

//HTTP HTTP
func HTTP(code int, err ...error) error {
	if len(err) > 0 {
		return &HTTPError{Code: code, error: err[0]}
	}

	return &HTTPError{Code: code, error: fmt.Errorf("%s", http.StatusText(code))}
}
