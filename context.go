package gen

import (
	"encoding/json"
	"encoding/xml"
	"net"
	"net/http"
	"strings"

	"github.com/x-io/gen/binding"
	"github.com/x-io/gen/core"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// Context defines request and response Context
type context struct {
	index  int
	router core.Routes
	route  core.Route
	params core.Params
	//query  core.Values

	data     map[string]interface{}
	request  *http.Request
	response *Response
	result   interface{}

	matched bool
	level   bool

	action interface{}
}

func newContext(router core.Routes) *context {
	c := new(context)
	c.router = router
	c.response = new(Response)
	return c
}

func (ctx *context) reset(w http.ResponseWriter, req *http.Request) {
	ctx.index = 0
	ctx.response.reset(w)
	ctx.request = req
	ctx.result = nil
	ctx.route = nil

	//ctx.query = nil
	ctx.params = nil
	ctx.action = nil
	ctx.level = true
	ctx.matched = false
	ctx.data = make(map[string]interface{})

}

func (ctx *context) newAction() {
	if !ctx.matched {
		if ctx.route == nil {
			ctx.route, ctx.params = ctx.router.Match(ctx.request.Method, ctx.request.URL.Path)
		}
		ctx.matched = true
	}
}

func (ctx *context) execute() {
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
		case func(core.Context):
			fn(ctx)
		}

		// not route matched
	} else if !ctx.response.Written() {
		//ctx.NotFound()
	}
}

func (ctx *context) invoke() {
	if ctx.level {
		if !ctx.router.Middleware(ctx, ctx.index) {
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
func (ctx *context) Next() {
	ctx.index++
	ctx.invoke()
}

// Req returns current HTTP Request information
func (ctx *context) Request() *http.Request {
	return ctx.request
}

func (ctx *context) Response() core.Response {
	return ctx.response
}

// Route returns route
func (ctx *context) Route() core.Route {
	ctx.newAction()
	return ctx.route
}

func (ctx *context) Result() interface{} {
	return ctx.result
}

// Params returns the URL params
func (ctx *context) Params() core.Params {
	ctx.newAction()
	return ctx.params
}

func (ctx *context) Header(name string) string {
	return ctx.request.Header.Get(name)
}

func (ctx *context) Query(name string) string {
	return ctx.request.FormValue(name)
}

func (ctx *context) Data(name string) string {
	return ctx.data[name].(string)
}

func (ctx *context) SetData(key string, value interface{}) {
	ctx.data[key] = value
}

func (ctx *context) Get(name string) core.Values {
	//fmt.Println(ctx.request.Header, ctx.data)
	if v, ok := ctx.request.Header[name]; ok {
		return core.Values(v[0])
	}

	if v, ok := ctx.data[name]; ok {
		return core.Values(v.(string))
	}

	if v := ctx.Query(name); v != "" {
		return core.Values(v)
	}

	return ""
}

// func (ctx *context) FormValue(name string) url.Values {
// 	return ctx.request.Form
// }

// Bind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
// 		"application/json" --> JSON binding
// 		"application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (ctx *context) Bind(obj interface{}) error {
	b := binding.Default(ctx.request.Method, ctx.ContentType())
	return ctx.BindWith(obj, b)
}

// BindJSON is a shortcut for c.BindWith(obj, binding.JSON)
func (ctx *context) BindJSON(obj interface{}) error {
	return ctx.BindWith(obj, binding.JSON)
}

// BindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (ctx *context) BindWith(obj interface{}, b binding.Binding) error {
	if err := b.Bind(ctx.request, obj); err != nil {
		ctx.Abort(400, err.Error()) //.SetType(ErrorTypeBind)
		return err
	}
	return nil
}

// IsAjax returns if the request is an ajax request
func (ctx *context) IsAjax() bool {
	return ctx.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (ctx *context) ClientIP() string {
	ips := strings.TrimSpace(ctx.request.Header.Get("X-Real-Ip"))
	if len(ips) > 0 {
		return ips
	}

	ips = ctx.request.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(ips, ','); index >= 0 {
		ips = ips[0:index]
	}
	ips = strings.TrimSpace(ips)
	if len(ips) > 0 {
		return ips
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.request.RemoteAddr)); err == nil {
		return ip
	}
	return "127.0.0.1"
}

// ContentType returns the Content-Type header of the request.
func (ctx *context) ContentType() string {
	return ctx.request.Header.Get("Content-Type")
}

// func (ctx *context) SetHeader(key, value string) {
// 	ctx.request.Header.Set(key, value)
// }

// // SecureCookies generates a secret cookie
// func (ctx *context) SecureCookies(secret string) Cookies {
// 	return &secureCookies{
// 		(*cookies)(ctx),
// 		secret,
// 	}
// }

// // Cookies returns the cookies
// func (ctx *context) Cookies() Cookies {
// 	return (*cookies)(ctx)
// }

// // Forms returns the query names and values
// func (ctx *context) Forms() *Forms {
// 	ctx.req.ParseForm()
// 	return (*Forms)(ctx.req)
// }

// Action returns action
func (ctx *context) Action() interface{} {
	ctx.newAction()
	return ctx.action
}

// // ActionValue returns action value
// func (ctx *context) ActionValue() reflect.Value {
// 	ctx.newAction()
// 	return ctx.callArgs[0]
// }

// // ActionTag returns field tag on action struct
// // TODO: cache the name
// func (ctx *context) ActionTag(fieldName string) string {
// 	ctx.newAction()
// 	if ctx.route.routeType == StructPtrRoute || ctx.route.routeType == StructRoute {
// 		tp := ctx.callArgs[0].Type()
// 		if tp.Kind() == reflect.Ptr {
// 			tp = tp.Elem()
// 		}
// 		field, ok := tp.FieldByName(fieldName)
// 		if !ok {
// 			return ""
// 		}
// 		return string(field.Tag)
// 	}
// 	return ""
// }

// // WriteString writes a string to response write
// func (ctx *context) WriteString(content string) (int, error) {
// 	return io.WriteString(ctx.ResponseWriter, content)
// }

// Redirect redirects the request to another URL
func (ctx *context) Redirect(url string, status ...int) {
	s := http.StatusFound
	if len(status) > 0 {
		s = status[0]
	}
	http.Redirect(ctx.response, ctx.request, url, s)
}

// NoContent writes a 204 HTTP response
func (ctx *context) NoContent(message ...string) {
	ctx.Abort(http.StatusNoContent, message...)
}

// NotModified writes a 304 HTTP response
func (ctx *context) NotModified(message ...string) {
	ctx.Abort(http.StatusNotModified, message...)
}

// BadRequest writes a 400 HTTP response
func (ctx *context) BadRequest(message ...string) {
	ctx.Abort(http.StatusBadRequest, message...)
}

// Unauthorized writes a 401 HTTP response
func (ctx *context) Unauthorized(message ...string) {
	ctx.Abort(http.StatusUnauthorized, message...)
}

// NotFound writes a 404 HTTP response
func (ctx *context) NotFound(message ...string) {
	ctx.Abort(http.StatusNotFound, message...)
}

func (ctx *context) Error(status int, body ...error) {
	ctx.result = body[0]
	ctx.response.WriteHeader(status)
	if len(body) == 0 {
		ctx.response.WriteString(http.StatusText(status))
	} else {
		ctx.response.WriteString(body[0].Error())
	}
	ctx.index = 100
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (ctx *context) Abort(status int, body ...string) {
	if len(body) > 0 {
		ctx.result = body[0]
	} else {
		ctx.result = ""
	}

	ctx.response.WriteHeader(status)
	if len(body) == 0 {
		ctx.response.WriteString(http.StatusText(status))
	} else {
		ctx.response.WriteString(body[0])
	}
	ctx.index = 100
}

//AbortJSON AbortJSON
func (ctx *context) AbortJSON(status int, obj interface{}) {
	ctx.result = "obj"
	ctx.response.WriteHeader(status)
	ctx.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(ctx.response).Encode(obj)
	if err != nil {
		ctx.response.Header().Del("Content-Type")
	}
	ctx.index = 100
}

//ToJSON serves marshaled JSON content from obj
func (ctx *context) ToJSON(obj interface{}) error {
	ctx.result = "obj"
	ctx.response.WriteHeader(200)
	ctx.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(ctx.response).Encode(obj)
	if err != nil {
		ctx.response.Header().Del("Content-Type")
	}
	return err
}

//ToXML serves marshaled XML content from obj
func (ctx *context) ToXML(obj interface{}) error {
	ctx.result = "obj"
	ctx.response.WriteHeader(200)
	ctx.response.Header().Set("Content-Type", "application/xml; charset=utf-8")
	err := xml.NewEncoder(ctx.response).Encode(obj)
	if err != nil {
		ctx.response.Header().Del("Content-Type")
	}
	return err
}

func (ctx *context) ToString(obj string) error {
	ctx.result = "obj"
	ctx.response.WriteHeader(200)
	ctx.response.Header().Set("Content-Type", "application/text; charset=utf-8")
	_, err := ctx.response.WriteString(obj)
	if err != nil {
		ctx.response.Header().Del("Content-Type")
	}
	return err
}