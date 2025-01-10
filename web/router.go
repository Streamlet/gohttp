package web

import (
	"net/http"
)

type Router struct {
	http.ServeMux
	cacheProvider CacheProvider
	customContext interface{}
}

func NewRouter(cacheProvider CacheProvider, customContext interface{}) *Router {
	return &Router{cacheProvider: cacheProvider, customContext: customContext}
}

func wrapHandler(handler func(Context), cacheProvider CacheProvider, customContext interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newHttpContext(w, r, cacheProvider, customContext)
		defer ctx.Release()
		defer func() {
			if err := recover(); err != nil {
				ctx.HttpError(http.StatusInternalServerError)
			}
		}()
		handler(ctx)
	}
}

func (r *Router) Handle(pattern string, handler func(Context)) {
	r.ServeMux.HandleFunc(pattern, wrapHandler(handler, r.cacheProvider, r.customContext))
}
