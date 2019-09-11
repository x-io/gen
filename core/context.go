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

	// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
	BindJSON(obj interface{}) error

	// // BindWith binds the passed struct pointer using the specified binding engine.
	// // See the binding package.
	// BindWith(obj interface{}, b binding.Binding) error
	Action() interface{}

	Result() interface{}
	ClientIP() string
	IsAjax() bool

	NoContent(message ...string)

	// NotModified writes a 304 HTTP response
	NotModified(message ...string)

	// Unauthorized writes a 401 HTTP response
	// Unauthorized(message ...string)

	// NotFound writes a 404 HTTP response
	NotFound(message ...string)
	// Abort is a helper method that sends an HTTP header and an optional
	// body. It is useful for returning 4xx or 5xx errors.
	// Once it has been called, any return value from the handler will
	// not be written to the response.
	Abort(status int, body ...string)
	AbortJSON(status int, obj interface{})

	Error(status int, err ...error)
	
	Redirect(url string, status ...int)
	ToJSON(obj interface{}) error
	ToXML(obj interface{}) error
	ToString(obj string) error
}
