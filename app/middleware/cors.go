package middleware

import (
	. "github.com/milkbobo/fishgoweb/app/cors"
	. "github.com/milkbobo/fishgoweb/app/router"
	"net/http"
)

func NewCorsMiddleware(cors Cors) RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				cors.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
					last(w, r, param)
				})
			},
		}
	}
}
