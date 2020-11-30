package middleware

import (
	. "github.com/milkbobo/fishgoweb/app/gzip"
	. "github.com/milkbobo/fishgoweb/app/router"
	"net/http"
)

func NewGzipMiddleware(gzip Gzip) RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				gzip.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
					last(w, r, param)
				})
			},
		}
	}
}
