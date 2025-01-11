package gohttp

import (
	"net/http"
)

type RouterInterface[T HttpContext] interface {
	http.Handler
	Handle(pattern string, handler func(T))
}

func NewRouter[T HttpContext](cp CacheProvider, cf ContextFactory[T]) RouterInterface[T] {
	if cp == nil {
		cp = newMemoryCache()
	}
	return &router[T]{sessionManager: newSessionManager(cp), contextFactory: cf}
}

type router[T HttpContext] struct {
	serveMux       http.ServeMux
	sessionManager SessionManager
	contextFactory ContextFactory[T]
}

type ContextFactory[T HttpContext] func(w http.ResponseWriter, r *http.Request, sm SessionManager) T

func (rt *router[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.serveMux.ServeHTTP(w, r)
}
func (rt *router[T]) Handle(pattern string, handler func(T)) {
	rt.serveMux.HandleFunc(pattern, wrapHandler(handler, rt.sessionManager, rt.contextFactory))
}

func wrapHandler[T HttpContext](handler func(T), sm SessionManager, cf ContextFactory[T]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := cf(w, r, sm)
		defer func() {
			if err := recover(); err != nil {
				ctx.HttpError(http.StatusInternalServerError)
			}
		}()
		handler(ctx)
	}
}
