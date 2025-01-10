package web

import (
	"net/http"
)

type Router[T any] struct {
	http.ServeMux
	cacheProvider CacheProvider
	extContext    T
}

func NewRouter[T any](cacheProvider CacheProvider, extContext T) *Router[T] {
	return &Router[T]{cacheProvider: cacheProvider, extContext: extContext}
}

func wrapHandler[T any](handler func(Context[T]), cacheProvider CacheProvider, extContext T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newHttpContext[T](w, r, cacheProvider, extContext)
		defer ctx.Release()
		defer func() {
			if err := recover(); err != nil {
				ctx.HttpError(http.StatusInternalServerError)
			}
		}()
		handler(ctx)
	}
}

func (r *Router[T]) Handle(pattern string, handler func(Context[T])) {
	r.ServeMux.HandleFunc(pattern, wrapHandler(handler, r.cacheProvider, r.extContext))
}
