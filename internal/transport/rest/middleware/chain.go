package middleware

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func Chain(mux http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		mux = middlewares[i](mux)
	}
	return mux
}
