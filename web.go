package gen

import (
	"github.com/x-io/gen/core"
	"github.com/x-io/gen/middlewares/io"
	"github.com/x-io/gen/middlewares/logger"
	"github.com/x-io/gen/middlewares/statics"
)

//D Data
type D map[string]interface{}

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
		io.Middleware(false),
		statics.Middleware(statics.Options{H5History: false}),
	)
}

// H5History returns an Engine instance with the Logger and Recovery middleware already attached.
func H5History() *Server {
	return New(
		logger.Middleware(),
		io.Middleware(false),
		// statics.Middleware(statics.Config{Prefix: "static"}),
		statics.Middleware(statics.Options{H5History: true}),
	)
}
