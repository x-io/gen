package io

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/x-io/gen/binding"
	"github.com/x-io/gen/core"
	"github.com/x-io/gen/errors"
)

//Middleware Middleware
func Middleware(debug bool) core.Middleware {
	out := log.Writer()

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
				ctx.Write(errors.HTTP(500))

				//	ctx.Abort(500)

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
		response := ctx.Response()
		if response.Written() {
			return
		}

		var content interface{}

		switch res := ctx.Result().(type) {
		case core.HTTPError:
			if res.Status() > 0 {
				response.WriteHeader(res.Status())
				response.WriteString(http.StatusText(res.Status()))
			} else {
				response.WriteHeader(http.StatusBadRequest)
				response.WriteString(res.Error())
			}
		case core.Error:
			response.WriteHeader(http.StatusBadRequest)
			content = res
		case error:
			if res == sql.ErrNoRows {
				response.WriteHeader(http.StatusNoContent)
			} else {
				response.WriteHeader(http.StatusBadRequest)
				response.WriteString(res.Error())
			}
		case string:
			response.WriteHeader(http.StatusOK)
			response.WriteString(res)
		case []byte:
			response.WriteHeader(http.StatusOK)
			response.Write(res)
		default:
			response.WriteHeader(http.StatusOK)
			content = res
		}

		if content != nil {

			b := binding.GetBinding(ctx.Request().Method, ctx.ContentType())

			if err := b.Write(response, content); err != nil {
				//return err
			}
		}
	}
}
