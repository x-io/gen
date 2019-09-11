package recovery

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"os"
	"runtime"
	"time"

	"github.com/x-io/gen/core"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func Middleware(debug bool) core.Middleware {
	return LoggerWithWriter(os.Stdout, debug)
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, debug bool) core.Middleware {
	return func(ctx core.Context) {
		defer func() {
			if err := recover(); err != nil {
				if debug {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(ctx.Request(), true)
					fmt.Fprintf(out, "[Gen] %v [Recovery] panic recovered: %s\n%s\n%s\n",
						time.Now().Format("2006/01/02 - 15:04:05"),
						err,
						string(httprequest),
						stack,
					)
				} else {
					fmt.Fprintf(out, "[Gen] %v [Recovery] panic recovered: %s\n",
						time.Now().Format("2006/01/02 - 15:04:05"),
						err,
					)
				}

				ctx.Abort(500)

				// var buf bytes.Buffer
				// fmt.Fprintf(&buf, "Handler crashed with error: %v", err)

				// for i := 1; ; i++ {
				// 	_, file, line, ok := runtime.Caller(i)
				// 	if !ok {
				// 		break
				// 	} else {
				// 		fmt.Fprintf(&buf, "\n")
				// 	}
				// 	fmt.Fprintf(&buf, "%v:%v", file, line)
				// }
				//log.Println(e)
				// var content = buf.String()
				// ctx.Logger.Error(content)

				// if !ctx.Written() {
				// 	if !debug {
				// 		ctx.Result = InternalServerError(http.StatusText(http.StatusInternalServerError))
				// 	} else {
				// 		ctx.Result = InternalServerError(content)
				// 	}
				// }
			}
		}()

		ctx.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
