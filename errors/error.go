package errors

import (
	"fmt"
	"net/http"
)

// Error Error
type Error struct {
	error
	HttpCode int
	Message  string `json:"message,omitempty"`
}

// Error Error
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.error.Error()
}

// HttpStatus HttpStatus
func (e *Error) HttpStatus() int {
	return e.HttpCode
}

// New New
func New(msg string) error {
	return &Error{Message: msg}
}

// HTTP HTTP
func HTTP(code int, err ...error) error {
	if len(err) > 0 {
		return &Error{HttpCode: code, error: err[0]}
	}

	return &Error{HttpCode: code, error: fmt.Errorf("%s", http.StatusText(code))}
}
