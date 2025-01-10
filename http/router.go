package http

import (
	"net/http"
)

type Router struct {
	http.ServeMux
	cacheProvider CacheProvider
}

func NewRouter(cacheProvider CacheProvider) *Router {
	return &Router{cacheProvider: cacheProvider}
}

func wrapHandler(handler func(Context), cacheProvider CacheProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newHttpContext(w, r, cacheProvider)
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
	r.ServeMux.HandleFunc(pattern, wrapHandler(handler, r.cacheProvider))
}
