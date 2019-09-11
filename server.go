package gen

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/x-io/gen/core"
	"github.com/x-io/gen/log"
	"github.com/x-io/gen/router"
)

type (
	//Server Web Server
	Server struct {
		http.Server
		core.Router
		pool sync.Pool
	}
)

func newServer() *Server {
	e := new(Server)
	e.Handler = e
	route := router.New()
	e.Router = route
	e.pool.New = func() interface{} {
		return newContext(route)
	}
	return e
}

// ServeHTTP implementes net/http interface so that it could run with net/http
func (e *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := e.pool.Get().(*context)
	c.reset(w, req)
	c.invoke()
	e.pool.Put(c)
}

// Run the http server. Listening on os.GetEnv("PORT") or 8000 by default.
func (e *Server) Run(args ...interface{}) (err error) {
	addr := getAddress(args...)
	log.Info("Listening on http://" + addr)

	// if e.shuttingDown() {
	// 	return ErrServerClosed
	// }

	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	//defer ln.Close()
	go func() {
		e.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
		ln.Close()
	}()

	return
}

// RunTLS runs the https server with special cert and key files
func (e *Server) RunTLS(certFile, keyFile string, args ...interface{}) (err error) {
	addr := getAddress(args...)
	log.Info("Listening on https://" + addr)
	e.Addr = addr
	err = e.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Error(err)
	}

	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	//defer ln.Close()
	go func() {
		e.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, certFile, keyFile)
		ln.Close()
	}()

	return
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
