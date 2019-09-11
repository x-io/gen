package error

import (
	"fmt"
	"net/http"

	"github.com/x-io/gen/core"
)

// AbortError defines an interface to describe HTTP error
type AbortError interface {
	error
	Code() int
}

type abortError struct {
	code    int
	content string
}

func (a *abortError) Code() int {
	return a.code
}

func (a *abortError) Error() string {
	return fmt.Sprintf("%v", a.content)
}

// Abort returns an AbortError
func Abort(code int, content ...string) AbortError {
	if len(content) >= 1 {
		return &abortError{code, content[0]}
	}
	return &abortError{code, http.StatusText(code)}
}

// NotFound returns not found HTTP error
func NotFound(content ...string) AbortError {
	return Abort(http.StatusNotFound, content...)
}

// NotSupported returns not supported HTTP error
func NotSupported(content ...string) AbortError {
	return Abort(http.StatusMethodNotAllowed, content...)
}

// InternalServerError returns internal server HTTP error
func InternalServerError(content ...string) AbortError {
	return Abort(http.StatusInternalServerError, content...)
}

// Forbidden returns forbidden HTTP error
func Forbidden(content ...string) AbortError {
	return Abort(http.StatusForbidden, content...)
}

// Unauthorized returns unauthorized HTTP error
func Unauthorized(content ...string) AbortError {
	return Abort(http.StatusUnauthorized, content...)
}

// Errors returns default errorhandler, you can use your self handler
func Errors() core.Middleware {
	return func(ctx core.Context) {
		ctx.Next()

		switch res := ctx.Result().(type) {
		case AbortError:
			ctx.Response().WriteHeader(res.Code())
			ctx.Response().WriteString(res.Error())
		case error:
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			ctx.Response().WriteString(res.Error())
		default:
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			ctx.Response().WriteString(http.StatusText(http.StatusInternalServerError))
		}
	}
}
