package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/x-io/gen/core"
	"github.com/x-io/gen/errors"
)

var (
	defaultBearer = "Bearer"
	defaultKey    = "X-IO"
)

//Options Options
type Options struct {
	Key            string
	Bearer         string
	CheckWebSocket bool
	Exclude        []string
}

//NewOption NewOption
func NewOption() *Options {
	return &Options{}
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.Bearer) == 0 {
		opt.Bearer = defaultBearer
	}
	if len(opt.Key) == 0 {
		opt.Key = defaultKey
	}

	return opt
}

// Middleware new create a JWT Middleware
func Middleware(options ...Options) core.Middleware {
	option := prepareOptions(options)
	return func(ctx core.Context) {
		request := ctx.Request()

		if option.CheckWebSocket {
			// Skip WebSocket
			if (request.Header.Get("Upgrade")) == "WebSocket" {
				ctx.Next()
				return
			}
		}

		auth := request.Header.Get("Authorization")
		l := len(option.Bearer)
		if len(auth) > l+1 && auth[:l] == option.Bearer {
			token, err := jwt.Parse(auth[l+1:], func(token *jwt.Token) (interface{}, error) {
				// Always check the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				// Return the key for validation
				return []byte(option.Key), nil
			})

			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

					for k, v := range claims {
						ctx.SetData(k, v)
					}

					ctx.Next()
					return
				}
			}

		}

		if !isContain(option.Exclude, request.URL.Path) {
			ctx.Write(errors.HTTP(http.StatusUnauthorized))
			return
		}

		ctx.Next()
	}
}

func isContain(items []string, item string) bool {
	for _, v := range items {
		if v == "*" || strings.Contains(item, v) {
			return true
		}
	}
	return false
}
