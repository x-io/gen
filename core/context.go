package core

import (
	"net/http"
)

// Context defines context interface
type Context interface {
	Next()
	Request() *http.Request
	Response() Response
	Route() Route
	Params() Params

	Get(name string) Values
	Header(name string) string
	Query(name string) string
	Data(name string) string
	SetData(key string, value interface{})

	Bind(obj interface{}) error
	Write(obj interface{}) error

	// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
	BindJSON(obj interface{}) error

	// BindWith binds the passed struct pointer using the specified binding engine.
	// See the binding package.
	//BindWith(obj interface{}, b binding.Binding) error
	Action() interface{}

	Result() interface{}
	ClientIP() string
	IsAjax() bool

	// NoContent writes a 204 HTTP response
	NoContent(message ...interface{})

	// NotModified writes a 304 HTTP response
	NotModified(message ...interface{})

	// BadRequest writes a 400 HTTP response
	BadRequest(message ...interface{})

	// Unauthorized writes a 401 HTTP response
	Unauthorized(message ...interface{})

	// NotFound writes a 404 HTTP response
	NotFound(message ...interface{})

	// Abort is a helper method that sends an HTTP header and an optional
	// body. It is useful for returning 4xx or 5xx errors.
	// Once it has been called, any return value from the handler will
	// not be written to the response.
	Abort(status int, body ...interface{})

	Error(status int, err ...error)

	Redirect(url string, status ...int)
	ToJSON(obj interface{}) error
	ToXML(obj interface{}) error
	ToString(obj string) error
}
