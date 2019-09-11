package core

// Middleware defines middleware interface
// Middleware describes the handle function
type Middleware func(Context)

// Handle executes the handler
func (h Middleware) Handle(ctx Context) {
	h(ctx)
}
