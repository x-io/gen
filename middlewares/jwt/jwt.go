package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/x-io/gen/core"
)

var (
	defaultBearer = "Bearer"
	defaultKey    = "JWT"
	claimKey      = ""
	bearerKey     = ""
)

type auther interface {
	SetClaims(map[string]interface{})
	GetClaim(string) interface{}
}

type Auther map[string]interface{}

func (a Auther) SetClaims(claims map[string]interface{}) {
	a = claims
}

func (a Auther) GetClaim(key string) interface{} {
	return a[key]
}

//Options Options
type Options struct {
	Bearer         string
	Key            string
	CheckWebSocket bool
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

	// Defaults
	if len(opt.Key) == 0 {
		opt.Key = defaultKey
	}

	if len(opt.Bearer) == 0 {
		opt.Bearer = defaultBearer
	}
	return opt
}

// Middleware new create a JWT Middleware
func Middleware(options ...Options) core.Middleware {
	option := prepareOptions(options)
	claimKey = option.Key
	bearerKey = option.Bearer
	return func(ctx core.Context) {
		// response := ctx.Response()

		request := ctx.Request()

		if option.CheckWebSocket {
			// Skip WebSocket
			if (request.Header.Get("Upgrade")) == "WebSocket" {
				ctx.Next()
				return
			}
		}

		//if a, ok := ctx.Action().(auther); ok {

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
			ctx.Abort(http.StatusForbidden)
			return
		}
		// ctx.Unauthorized()
		ctx.Abort(http.StatusUnauthorized)
		return
		//}

		//ctx.Next()
	}
}

//NewToken NewToken
func NewToken(claims ...map[string]interface{}) (string, error) {
	claim := make(jwt.MapClaims)
	claim["exp"] = time.Now().Add(time.Hour * 72).Unix()
	if len(claims) > 0 {
		for k, v := range claims[0] {
			claim[k] = v
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(claimKey))
}

//Parse Parse
func Parse(tokenString string) (jwt.MapClaims, error) {
	l := len(bearerKey)
	token, err := jwt.Parse(tokenString[l+1:], func(token *jwt.Token) (interface{}, error) {
		// Always check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the key for validation
		return []byte(claimKey), nil
	})

	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			return claims, nil
		}
	}

	return nil, err
}

// //NewToken NewToken
// func NewToken(key string, claims ...map[string]interface{}) (string, error) {
// 	claim := make(jwt.MapClaims)
// 	claim["exp"] = time.Now().Add(time.Hour * 72).Unix()
// 	if len(claims) > 0 {
// 		for k, v := range claims[0] {
// 			claim[k] = v
// 		}
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

// 	return token.SignedString([]byte(key))
// }
