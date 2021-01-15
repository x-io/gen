package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/x-io/gen/core"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})

	out   io.Writer = os.Stderr
	color           = true
)

// Middleware instances a Logger middleware that will write the logs to gen.DefaultWriter
// By default gen.DefaultWriter = os.Stdout
func Middleware() core.Middleware {
	return WithWriter()
}

//SetWriter SetWriter
func SetWriter(w io.Writer) {
	out = w
	color = w == os.Stderr
}

// WithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func WithWriter(notlogged ...string) core.Middleware {
	isWindows := false
	if runtime.GOOS == "windows" {
		isWindows = true
	}
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *core.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		path := c.Request.URL.Path
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)
			method := c.Request.Method

			clientIP := c.ClientIP()
			statusCode := c.Response.Status()

			path = c.Request.URL.String()

			//	comment := "" //c.Errors.ByType(ErrorTypePrivate).String()

			if statusCode > 200 {
				if result := c.Result; result != nil {
					fmt.Fprintf(out, "[Gen] %s \n", result)
				}
			}

			if !color || isWindows {
				fmt.Fprintf(out, "[Gen] %v |%3d| %9v | %s | %-7s %s\n",
					end.Format("2006/01/02 - 15:04:05"),
					statusCode,
					latency,
					clientIP,
					method,
					path,
					//comment,
				)
			} else {
				fmt.Fprintf(out, "[Gen] %v |%s %3d %s| %9v | %s |%s %s %-7s %s\n",
					end.Format("2006/01/02 - 15:04:05"),
					colorForStatus(statusCode), statusCode, reset,
					latency,
					clientIP,
					colorForMethod(method), reset, method,
					path,
					//comment,
				)
			}
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
