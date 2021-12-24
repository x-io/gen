package router

import (
	"regexp"
	"strings"

	"github.com/x-io/gen/core"
)

type routeTable map[string]*node

//Router Router
type Router struct {
	level       int
	prefix      string
	table       routeTable
	middlewares []core.Middleware
	meta        map[string]string
}

//New NewRouter
func New() *Router {
	r := new(Router)

	r.level = 0
	r.prefix = ""
	r.table = make(routeTable)
	r.middlewares = make([]core.Middleware, 0)
	return r
}

func newGroup(table routeTable, prefix string, level int, m ...core.Middleware) core.Router {
	r := new(Router)
	r.level = level
	r.prefix = prefix
	if len(m) > 0 {
		r.middlewares = m
	}
	r.table = table
	return r
}

func (r *Router) add(methods []string, path string, handler interface{}, m ...core.Middleware) core.Route {
	if r.level == 0 {
		return r.Route(methods, r.prefix+path, handler, m...)
	} else {
		return r.Route(methods, r.prefix+path, handler, r.middlewares...).Use(m...)
	}
}

//Use Use
func (r *Router) Use(m ...core.Middleware) core.Router {
	if r.middlewares == nil {
		r.middlewares = make([]core.Middleware, 0)
	}

	r.middlewares = append(r.middlewares, m...)
	return r
}

//Use Meta
func (r *Router) Meta(name, value string) core.Router {
	if r.meta == nil {
		r.meta = make(map[string]string)
	}

	r.meta[name] = value
	return r
}

//Group Group
func (r *Router) Group(prefix string, m ...core.Middleware) core.Router {
	if r.level == 0 {
		return newGroup(r.table, r.prefix+prefix, r.level+1, m...)
	}

	return newGroup(r.table, r.prefix+prefix, r.level+1, r.middlewares...).Use(m...)
}

// Get sets a route with GET method
func (r *Router) Get(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"GET"}, path, handler, m...)
}

// Post sets a route with POST method
func (r *Router) Post(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"POST"}, path, handler, m...)
}

// Put sets a route with PUT method
func (r *Router) Put(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"PUT"}, path, handler, m...)
}

// Delete sets a route with DELETE method
func (r *Router) Delete(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"DELETE"}, path, handler, m...)
}

// Head sets a route with HEAD method
func (r *Router) Head(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"HEAD"}, path, handler, m...)
}

// Options sets a route with OPTIONS method
func (r *Router) Options(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"OPTIONS"}, path, handler, m...)
}

// Trace sets a route with TRACE method
func (r *Router) Trace(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"TRACE"}, path, handler, m...)
}

// Patch sets a route with PATCH method
func (r *Router) Patch(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"PATCH"}, path, handler, m...)
}

// Any sets a route every support method is OK.
func (r *Router) Any(path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "TRACE", "PATCH"}, path, handler, m...)
}

//Handle Handle
func (r *Router) Handle(method string, path string, handler interface{}, m ...core.Middleware) core.Route {
	return r.add([]string{method}, path, handler, m...)
}

// Route adds route
func (r *Router) Route(methods []string, path string, handler interface{}, m ...core.Middleware) core.Route {

	url, meta := pathDecode(path)

	if url[0] != '/' {
		panic("path must begin with '/' in path '" + url + "'")
	}

	metas := make(map[string]string)

	if r.meta != nil {
		for n, v := range r.meta {
			metas[n] = v
		}
	}

	if len(meta) > 0 {
		if v, ok := metas["route"]; ok {
			metas["route"] = v + "," + meta
		} else {
			metas["route"] = meta
		}
	}

	route := &Route{
		handle:      handler,
		meta:        metas,
		middlewares: m,
	}

	for _, v := range methods {
		root := r.table[v]
		if root == nil {
			root = new(node)
			r.table[v] = root
		}

		root.addRoute(url, route)
	}

	return route
}

// Match for request path, match router
func (r *Router) Match(method, path string) (core.Route, core.Params) {
	path = strings.TrimSuffix(path, "/")
	if root := r.table[method]; root != nil {
		if route, params, _ := root.getValue(path); route != nil {
			return route, &params
		}
	}
	if root := r.table["*"]; root != nil {
		if route, params, _ := root.getValue(path); route != nil {
			return route, &params
		}
	}

	return nil, nil
}

//Middleware Middleware
func (r *Router) Middleware(ctx *core.Context, index int) bool {
	if index < len(r.middlewares) {
		r.middlewares[index].Handle(ctx)
		return true
	}
	return false
}

// Route defines HTTP route
type Route struct {
	handle      interface{}
	meta        map[string]string
	middlewares []core.Middleware
}

// Use addes some Route middlewares
func (e *Route) Use(m ...core.Middleware) core.Route {
	e.middlewares = append(e.middlewares, m...)
	return e
}

// Middleware Middleware
func (e *Route) Middleware(ctx *core.Context, index int) bool {
	if index < len(e.middlewares) {
		e.middlewares[index].Handle(ctx)
		return true
	}
	return false
}

//Handle Handle
func (e *Route) Handle() interface{} {
	return e.handle
}

func (e *Route) Meta(name string) string {
	return e.meta[name]
}

func pathDecode(data string) (string, string) {

	//解析正则表达式，如果成功返回解释器
	reg := regexp.MustCompile(`\[([\w\,\|]+)\]`)
	if reg == nil { //解释失败，返回nil
		return data, ""
	}

	path := data
	metas := make([]string, 0)
	matchs := reg.FindAllStringSubmatch(path, -1)

	for _, v := range matchs {
		if len(v) > 0 {
			metas = append(metas, v[1])
			path = strings.Replace(path, v[0], "", -1)
		}
	}
	return path, strings.Join(metas, ",")
}
