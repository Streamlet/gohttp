package main

import (
	"net/http"

	"github.com/Streamlet/gohttp"
)

type HttpContext interface {
	gohttp.HttpContext
	Custom() string
}

func NewContextFactory() gohttp.ContextFactory[HttpContext] {
	return &contextFactory{gohttp.NewSessionManager(nil)}
}

type contextFactory struct {
	sm gohttp.SessionManager
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &httpContext{
		gohttp.NewHttpContext(w, r, cf.sm),
	}
}

type httpContext struct {
	gohttp.HttpContext
}

func (cc *httpContext) Custom() string {
	return "custom"
}
