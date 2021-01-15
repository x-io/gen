package cors

import (
	"strconv"

	"github.com/x-io/gen/core"
)

//Options CROS配置
type Options struct {
	Enabled          bool
	AllowOrigin      string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
	AllowP3P         string
	MaxAge           int
}

//NewOptions NewOptions
func NewOptions() *Options {
	return &Options{}
}

//NewDefault NewDefault
func NewDefault() *Options {
	c := new(Options)
	c.Enabled = true
	c.AllowOrigin = "*"
	c.AllowMethods = "GET, POST, PUT, DELETE, OPTIONS"
	c.AllowHeaders = "Authorization, Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With"
	c.AllowP3P = "CP=\"CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR\""
	return c
}

// Middleware new create a CORS Middleware
func Middleware(options *Options) core.Middleware {

	return func(ctx *core.Context) {
		response := ctx.Response

		if options.Enabled {
			response.SetHeader("Access-Control-Allow-Origin", options.AllowOrigin)
			response.SetHeader("Access-Control-Allow-Methods", options.AllowMethods)
			response.SetHeader("Access-Control-Allow-Headers", options.AllowHeaders)
			response.SetHeader("Access-Control-Expose-Headers", options.ExposeHeaders)
			response.SetHeader("Access-Control-Allow-Credentials", strconv.FormatBool(options.AllowCredentials))
			response.SetHeader("Access-Control-Max-Age", strconv.Itoa(options.MaxAge))
			response.SetHeader("P3P", options.AllowP3P)
		}
		ctx.Next()
	}
}
