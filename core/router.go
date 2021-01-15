package core

//Router Router
type Router interface {
	//Use Use
	Use(m ...Middleware) Router
	//Group Group
	Group(prefix string, m ...Middleware) Router
	// Get sets a route with GET method
	Get(path string, handler interface{}, m ...Middleware)
	// Post sets a route with POST method
	Post(path string, handler interface{}, m ...Middleware)
	// Put sets a route with PUT method
	Put(path string, handler interface{}, m ...Middleware)
	// Delete sets a route with DELETE method
	Delete(path string, handler interface{}, m ...Middleware)
	// Head sets a route with HEAD method
	Head(path string, handler interface{}, m ...Middleware)
	// Options sets a route with OPTIONS method
	Options(path string, handler interface{}, m ...Middleware)
	// Trace sets a route with TRACE method
	Trace(path string, handler interface{}, m ...Middleware)

	// Patch sets a route with PATCH method
	Patch(path string, handler interface{}, m ...Middleware)
	// Any sets a route every support method is OK.
	Any(path string, handler interface{}, m ...Middleware)

	Handle(method string, path string, handler interface{}, m ...Middleware)
}

//Routes describes the interface of route
type Routes interface {
	//Route(methods []string, path string, handler interface{}, middleware ...Middleware) Route
	Match(method, path string) (Route, Params)

	//Use(middleware ...Middleware)
	Middleware(*Context, int) bool
}

// Route defines HTTP route
type Route interface {
	Use(middleware ...Middleware)
	Middleware(*Context, int) bool
	Handle() interface{}
}
