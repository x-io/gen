package core

import "net/http"

//Binding Binding
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
	Write(*Response, interface{}) error
}
