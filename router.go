package gohttp

import (
	"log"
	"net/http"
	"runtime/debug"
)

type RouterInterface[T HttpContext] interface {
	http.Handler
	RawHandle(pattern string, handler http.Handler)
	Handle(pattern string, handler func(T))
}

type ContextFactory[T HttpContext] interface {
	NewContext(w http.ResponseWriter, r *http.Request) T
}

func NewContextFactory(cp CacheProvider, cookieDomain string) ContextFactory[HttpContext] {
	if cp == nil {
		cp = newMemoryCache()
	}
	return &contextFactory{NewSessionManager(cp), cookieDomain}
}

type contextFactory struct {
	sessionManager SessionManager
	cookieDomain   string
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return NewHttpContext(w, r, cf.sessionManager, cf.cookieDomain)
}

func NewRouter[T HttpContext](cf ContextFactory[T]) RouterInterface[T] {
	return &router[T]{contextFactory: cf}
}

type router[T HttpContext] struct {
	serveMux       http.ServeMux
	contextFactory ContextFactory[T]
}

func (rt *router[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.serveMux.ServeHTTP(w, r)
}

func (rt *router[T]) RawHandle(pattern string, handler http.Handler) {
	rt.serveMux.Handle(pattern, handler)
}

func (rt *router[T]) Handle(pattern string, handler func(T)) {
	rt.serveMux.HandleFunc(pattern, wrapHandler(handler, rt.contextFactory))
}

func wrapHandler[T HttpContext](handler func(T), cf ContextFactory[T]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := cf.NewContext(w, r)
		defer func() {
			if err := recover(); err != nil {
				log.Println(err, string(debug.Stack()))
				ctx.HttpError(http.StatusInternalServerError)
			}
		}()
		handler(ctx)
	}
}
