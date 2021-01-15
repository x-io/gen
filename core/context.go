package core

import (
	"net"
	"net/http"
	"strings"

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

func (ctx *Context) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx.index = 0
	ctx.Response.reset(w)
	ctx.Request = req
	ctx.Result = nil
	ctx.route = nil

	//ctx.query = nil
	ctx.params = nil
	ctx.action = nil
	ctx.level = true
	ctx.matched = false
	ctx.data = make(map[string]interface{})
	ctx.invoke()

}

func (ctx *Context) newAction() {
	if !ctx.matched {
		if ctx.route == nil {
			ctx.route, ctx.params = ctx.Router.Match(ctx.Request.Method, ctx.Request.URL.Path)
		}
		ctx.matched = true
	}
}

func (ctx *Context) execute() {
	ctx.newAction()
	// route is matched
	if ctx.route != nil {
		if ctx.level {
			ctx.index = 0
			ctx.level = false
			ctx.invoke()
			return
		}

		switch fn := ctx.route.Handle().(type) {
		case func(*Context) interface{}:
			if data := fn(ctx); data != nil {
				ctx.Result = data
			}
		case func(*Context) (interface{}, error):
			if data, err := fn(ctx); err != nil {
				ctx.Result = err
			} else {
				ctx.Result = data
			}
		case func(*Context) error:
			if err := fn(ctx); err != nil {
				ctx.Result = err
			}
		default:

		}
		// not route matched
	} else if !ctx.Response.Written() {
		ctx.Result = errors.HTTP(http.StatusNotFound)
		//	ctx.NotFound()
	}
}

func (ctx *Context) invoke() {
	if ctx.level {
		if !ctx.Router.Middleware(ctx, ctx.index) {
			ctx.execute()
		}
	} else {
		if !ctx.route.Middleware(ctx, ctx.index) {
			ctx.execute()
		}
	}
}

// Next call next middleware or action
// WARNING: don't invoke this method on action
func (ctx *Context) Next() {
	ctx.index++
	ctx.invoke()
}

// // Route returns route
// func (ctx *Context) Route() Route {
// 	ctx.newAction()
// 	return ctx.route
// }

// Params returns the URL params
func (ctx *Context) Params(name string) string {
	ctx.newAction()
	return ctx.params.Get(name)
}

// Header returns header params
func (ctx *Context) Header(name string) string {
	return ctx.Request.Header.Get(name)
}

// Query returns params
func (ctx *Context) Query(name string) string {
	return ctx.Request.FormValue(name)
}

// Querys returns params
func (ctx *Context) Querys(name string) Values {
	return Values(ctx.Request.FormValue(name))
}

//Data Data
func (ctx *Context) Data(name string) interface{} {

	if v, ok := ctx.data[name]; ok {
		return v
	}

	return nil
}

//SetData SetData
func (ctx *Context) SetData(name string, value interface{}) {
	ctx.data[name] = value
}

// //Data Data
// func (ctx *Context) Data(key, name string) Values {

// 	if v, ok := ctx.data[key]; ok {
// 		if vv, ok := v.(map[string]interface{}); ok {
// 			if vvv, ok := vv[key]; ok {
// 				return Values(vvv.(string))
// 			}
// 		}
// 	}

// 	return ""
// }

// //Get ...
// func (ctx *Context) Get(name string) Values {

// 	if v, ok := ctx.Request.Header[name]; ok {
// 		return Values(v[0])
// 	}

// 	if v, ok := ctx.Data[name]; ok {
// 		return Values(v.(string))
// 	}

// 	if v := ctx.Request.FormValue(name); v != "" {
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
func (ctx *Context) Bind(obj interface{}) error {
	if err := ctx.Binding.Bind(ctx.Request, obj); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) Write(obj interface{}) error {
	ctx.Result = obj
	return nil
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (ctx *Context) ClientIP() string {
	ips := strings.TrimSpace(ctx.Request.Header.Get("X-Real-Ip"))
	if len(ips) > 0 {
		return ips
	}

	ips = ctx.Request.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(ips, ','); index >= 0 {
		ips = ips[0:index]
	}
	ips = strings.TrimSpace(ips)
	if len(ips) > 0 {
		return ips
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}
	return "127.0.0.1"
}

// ContentType returns the Content-Type header of the request.
func (ctx *Context) ContentType() string {
	return ctx.Request.Header.Get("Content-Type")
}

// IsAjax returns if the request is an ajax request
func (ctx *Context) IsAjax() bool {
	return ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// // Action returns action
// func (ctx *Context) Action() interface{} {
// 	ctx.newAction()
// 	return ctx.action
// }

// Redirect redirects the request to another URL
func (ctx *Context) Redirect(url string, status ...int) error {
	s := http.StatusFound
	if len(status) > 0 {
		s = status[0]
	}
	http.Redirect(ctx.Response, ctx.Request, url, s)
	return nil
}
