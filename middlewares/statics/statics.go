package statics

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/x-io/gen/core"
)

// Options defines Static middleware's options
type Options struct {
	RootPath   string
	Prefix     string
	IndexFiles []string
	ListDir    bool
	H5History  bool
	FilterExts []string
	// FileSystem is the interface for supporting any implmentation of file system.
	FileSystem http.FileSystem
}

// IsFilterExt decribes if rPath's ext match filter ext
func (s *Options) IsFilterExt(rPath string) bool {
	rext := path.Ext(rPath)
	for _, ext := range s.FilterExts {
		if rext == ext {
			return true
		}
	}
	return false
}

func prepareStaticOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.RootPath) == 0 {
		opt.RootPath = "./public"
	}

	if len(opt.Prefix) > 0 {
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
	}

	if len(opt.IndexFiles) == 0 {
		opt.IndexFiles = []string{"index.html", "index.htm"}
	}

	if opt.FileSystem == nil {
		ps, _ := filepath.Abs(opt.RootPath)
		opt.FileSystem = http.Dir(ps)
	}

	return opt
}

// Middleware return a Static for serving static files
func Middleware(options ...Options) core.Middleware {

	return func(ctx core.Context) {
		request := ctx.Request()
		if request.Method != "GET" && request.Method != "HEAD" {
			ctx.Next()
			return
		}

		rPath := request.URL.Path
		opt := prepareStaticOptions(options)

		// if defined prefix, then only check prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(rPath, opt.Prefix) {
				ctx.Next()
				return
			}

			if len(opt.Prefix) == len(rPath) {
				rPath = ""
			} else {
				rPath = rPath[len(opt.Prefix):]
			}
		}

		f, err := opt.FileSystem.Open(strings.TrimLeft(rPath, "/"))

		if err != nil {
			if os.IsNotExist(err) {
				if opt.Prefix != "" {
					ctx.NotFound()
				} else {
					ctx.Next()
					if ctx.Result() == nil {
						if opt.H5History {

							//fmt.Println("H5History", rPath)

							//try serving index.html or index.htm
							if len(opt.IndexFiles) > 0 {
								for _, index := range opt.IndexFiles {
									fi, err := opt.FileSystem.Open(strings.TrimLeft(index, "/"))
									if err != nil {
										if !os.IsNotExist(err) {
											ctx.Abort(http.StatusInternalServerError, err.Error())
											return
										}
									} else {
										defer fi.Close()
										finfo, err := fi.Stat()
										if err != nil {
											ctx.Abort(http.StatusInternalServerError, err.Error())
											return
										}
										if !finfo.IsDir() {
											http.ServeContent(ctx.Response(), request, finfo.Name(), finfo.ModTime(), fi)
											return
										}
									}
								}
							}
							ctx.NotFound()
						}
					}
				}
			} else {
				ctx.Abort(http.StatusInternalServerError, err.Error())
			}
			return
		}
		defer f.Close()

		finfo, err := f.Stat()
		if err != nil {
			ctx.Abort(http.StatusInternalServerError, err.Error())
			return
		}

		if !finfo.IsDir() {
			if len(opt.FilterExts) > 0 && !opt.IsFilterExt(rPath) {
				ctx.Next()
				return
			}

			http.ServeContent(ctx.Response(), request, finfo.Name(), finfo.ModTime(), f)
			return
		}

		// try serving index.html or index.htm
		if len(opt.IndexFiles) > 0 {
			for _, index := range opt.IndexFiles {
				fi, err := opt.FileSystem.Open(strings.TrimLeft(path.Join(rPath, index), "/"))
				if err != nil {
					if !os.IsNotExist(err) {
						ctx.Abort(http.StatusInternalServerError, err.Error())
						return
					}
				} else {
					finfo, err = fi.Stat()
					if err != nil {
						ctx.Abort(http.StatusInternalServerError, err.Error())
						return
					}
					if !finfo.IsDir() {
						http.ServeContent(ctx.Response(), request, finfo.Name(), finfo.ModTime(), fi)
						return
					}
				}
			}
		}
		ctx.Next()
	}
}
