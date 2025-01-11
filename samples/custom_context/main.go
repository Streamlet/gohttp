package main

import (
	"github.com/Streamlet/gohttp"
	"net/http"
)

type CustomContext interface {
	gohttp.HttpContext
	Custom() string
}

func NewContextFactory() gohttp.ContextFactory[CustomContext] {
	return &contextFactory{gohttp.NewSessionManager(nil)}
}

type contextFactory struct {
	sm gohttp.SessionManager
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) CustomContext {
	return &customContext{
		gohttp.NewHttpContext(w, r, cf.sm),
	}
}

type customContext struct {
	gohttp.HttpContext
}

func (cc *customContext) Custom() string {
	return "custom"
}

func ExtContextHandler(c CustomContext) {
	c.String(c.Custom())
}

func main() {
	application := gohttp.NewApplication[CustomContext](NewContextFactory())
	application.Handle("/", ExtContextHandler)
	application.ServePort(80)
}
