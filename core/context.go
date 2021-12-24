package core

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/x-io/gen/errors"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// Context defines request and response Context
type Context struct {
	index  int
	Router Routes
	route  Route
	params Params

	data map[string]interface{}

	Request  *http.Request
	Response *Response
	Result   interface{}

	Binding Binding

	matched bool
	level   bool

	action interface{}
}

func (c *Context) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.index = 0
	c.Response.reset(w)
	c.Request = req
	c.Result = nil
	c.route = nil

	//c.query = nil
	c.params = nil
	c.action = nil
	c.level = true
	c.matched = false
	c.data = make(map[string]interface{})
	c.invoke()

}

func (c *Context) newAction() {
	if !c.matched {
		if c.route == nil {
			c.route, c.params = c.Router.Match(c.Request.Method, c.Request.URL.Path)
		}
		c.matched = true
	}
}

func (c *Context) execute() {
	c.newAction()
	// route is matched
	if c.route != nil {
		if c.level {
			c.index = 0
			c.level = false
			c.invoke()
			return
		}

		switch fn := c.route.Handle().(type) {
		case func(*Context) interface{}:
			if data := fn(c); data != nil {
				c.Result = data
			}
		case func(*Context) (interface{}, error):
			if data, err := fn(c); err != nil {
				c.Result = err
			} else {
				c.Result = data
			}
		case func(*Context) error:
			if err := fn(c); err != nil {
				c.Result = err
			}
		default:

		}
		// not route matched
	} else if !c.Response.Written() {
		c.Result = errors.HTTP(http.StatusNotFound)
		//	c.NotFound()
	}
}

func (c *Context) invoke() {
	if c.level {
		if !c.Router.Middleware(c, c.index) {
			c.execute()
		}
	} else {
		if !c.route.Middleware(c, c.index) {
			c.execute()
		}
	}
}

// Next call next middleware or action
// WARNING: don't invoke this method on action
func (c *Context) Next() {
	c.index++
	c.invoke()
}

// // Route returns route
// func (c *Context) Route() Route {
// 	c.newAction()
// 	return c.route
// }

// Params returns the URL params
func (c *Context) Params(name string) string {
	c.newAction()
	return c.params.Get(name)
}

// Header returns header params
func (c *Context) Header(name string) string {
	return c.Request.Header.Get(name)
}

// Query returns params
func (c *Context) Query(name string) string {
	return c.Request.FormValue(name)
}

// Querys returns params
func (c *Context) Querys(name string) Values {
	return Values(c.Request.FormValue(name))
}

// Meta returns params
func (c *Context) Meta(name string) string {
	return c.route.Meta(name)
}

//Data Data
func (c *Context) Data(name string) interface{} {

	if v, ok := c.data[name]; ok {
		return v
	}

	return nil
}

//SetData SetData
func (c *Context) SetData(name string, value interface{}) {
	c.data[name] = value
}

// //Data Data
// func (c *Context) Data(key, name string) Values {

// 	if v, ok := c.data[key]; ok {
// 		if vv, ok := v.(map[string]interface{}); ok {
// 			if vvv, ok := vv[key]; ok {
// 				return Values(vvv.(string))
// 			}
// 		}
// 	}

// 	return ""
// }

// //Get ...
// func (c *Context) Get(name string) Values {

// 	if v, ok := c.Request.Header[name]; ok {
// 		return Values(v[0])
// 	}

// 	if v, ok := c.Data[name]; ok {
// 		return Values(v.(string))
// 	}

// 	if v := c.Request.FormValue(name); v != "" {
// 		return Values(v)
// 	}

// 	return ""
// }

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON binding
// 		"application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) Bind(obj interface{}) error {
	if err := c.Binding.Bind(c.Request, obj); err != nil {
		return err
	}
	return nil
}

func (c *Context) Write(obj interface{}) error {
	c.Result = obj
	return nil
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (c *Context) ClientIP() string {
	ips := strings.TrimSpace(c.Request.Header.Get("X-Real-Ip"))
	if len(ips) > 0 {
		return ips
	}

	ips = c.Request.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(ips, ','); index >= 0 {
		ips = ips[0:index]
	}
	ips = strings.TrimSpace(ips)
	if len(ips) > 0 {
		return ips
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return "127.0.0.1"
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return c.Request.Header.Get("Content-Type")
}

// IsAjax returns if the request is an ajax request
func (c *Context) IsAjax() bool {
	return c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// // Action returns action
// func (c *Context) Action() interface{} {
// 	c.newAction()
// 	return c.action
// }

// Redirect redirects the request to another URL
func (c *Context) Redirect(url string, status ...int) error {
	s := http.StatusFound
	if len(status) > 0 {
		s = status[0]
	}
	http.Redirect(c.Response, c.Request, url, s)
	return nil
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline always returns that there is no deadline (ok==false),
// maybe you want to use Request.Context().Deadline() instead.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to abort your work when the connection was closed
// you should use Request.Context().Done() instead.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Request.Context().Err() instead.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if name, ok := key.(string); ok {

		if v, ok := c.Request.Header[name]; ok {
			return v[0]
		}

		if v, ok := c.data[name]; ok {
			return v
		}

		if v := c.Request.FormValue(name); v != "" {
			return v
		}
	}
	return nil
}
