package core

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

//Response Response
type Response struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	status int
	size   int
}

func (rw *Response) reset(w http.ResponseWriter) {
	rw.ResponseWriter = w
	rw.status = 0
	rw.size = 0
}

//WriteHeader WriteHeader
func (rw *Response) WriteHeader(code int) {
	if code > 0 && rw.status != code {
		// if w.Written() {
		// 	debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		// }
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

//WriteHeaderNow WriteHeaderNow
func (rw *Response) WriteHeaderNow() {
	if !rw.Written() {
		rw.size = 0
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.ResponseWriter.WriteHeader(rw.status)
	}
}

func (rw *Response) Write(data []byte) (int, error) {
	rw.WriteHeaderNow()
	n, err := rw.ResponseWriter.Write(data)
	rw.size += n
	return n, err
}

//WriteString WriteString
func (rw *Response) WriteString(s string) (int, error) {
	rw.WriteHeaderNow()
	n, err := io.WriteString(rw.ResponseWriter, s)
	rw.size += n
	return n, err
}

//Status Status
func (rw *Response) Status() int {
	return rw.status
}

//Size Size
func (rw *Response) Size() int {
	return rw.size
}

//Written Written
func (rw *Response) Written() bool {
	return rw.status != 0
}

//Hijack Implements the http.Hijacker interface
func (rw *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if rw.size < 0 {
		rw.size = 0
	}
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

//CloseNotify Implements the http.CloseNotify interface
func (rw *Response) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

//Flush Implements the http.Flush interface
func (rw *Response) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}

//SetHeader SetHeader
func (rw *Response) SetHeader(key, value string) {
	rw.Header().Set(key, value)
}

//SetStatus SetStatus
func (rw *Response) SetStatus(code int) {
	if code > 0 && rw.status != code {
		// if w.Written() {
		// 	debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		// }
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}
