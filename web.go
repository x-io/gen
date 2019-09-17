package gen

import (
	"github.com/x-io/gen/core"
	"github.com/x-io/gen/middlewares/logger"
	"github.com/x-io/gen/middlewares/recovery"
	"github.com/x-io/gen/middlewares/statics"
)

//J Json
type J map[string]interface{}

// Version returns Framework's version
func Version() string {
	return "1.0.0"
}

// New creates tango with the default logger and handlers
func New(middlewares ...core.Middleware) *Server {
	e := newServer()
	e.Use(middlewares...)
	return e
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Server {
	return New(
		logger.Middleware(),
		recovery.Middleware(false),
		//	Compresses([]string{}),
		statics.Middleware(statics.Options{H5History: false}),
		//Return(),
		//Param(),
		//Contexts(),
	)
}

// H5History returns an Engine instance with the Logger and Recovery middleware already attached.
func H5History() *Server {
	return New(
		logger.Middleware(),
		recovery.Middleware(false),
		//	Compresses([]string{}),
		// statics.Middleware(statics.Config{Prefix: "static"}),
		statics.Middleware(statics.Options{H5History: true}),
		//Return(),
		//Param(),
		//Contexts(),
	)
}
