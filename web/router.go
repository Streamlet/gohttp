package web

import (
	"net/http"
)

type Router[T any] struct {
	http.ServeMux
	cacheProvider CacheProvider
	customContext T
}

func NewRouter[T any](cacheProvider CacheProvider, customContext T) *Router[T] {
	return &Router[T]{cacheProvider: cacheProvider, customContext: customContext}
}

func wrapHandler[T any](handler func(Context[T]), cacheProvider CacheProvider, customContext T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newHttpContext[T](w, r, cacheProvider, customContext)
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
	r.ServeMux.HandleFunc(pattern, wrapHandler(handler, r.cacheProvider, r.customContext))
}
