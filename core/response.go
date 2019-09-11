package core

import (
	"net/http"
)

// Response is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type Response interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Returns true if the response body was already written.
	Written() bool

	// Writes the string into the response body.
	WriteString(string) (int, error)
	// Forces to write the http header (status code + headers).
	//WriteHeaderNow()

	SetHeader(string, string)
	SetStatus(int)
}
