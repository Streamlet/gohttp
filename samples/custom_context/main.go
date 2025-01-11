package main

import (
	"github.com/Streamlet/gohttp"
	"net/http"
)

type CustomContext interface {
	gohttp.HttpContext
	Custom() string
}

func NewContext(w http.ResponseWriter, r *http.Request, sm gohttp.SessionManager) CustomContext {
	return &customContext{
		gohttp.NewHttpContext(w, r, sm),
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
	application := gohttp.NewApplication[CustomContext](NewContext, nil)
	application.Handle("/", ExtContextHandler)
	application.ServePort(80)
}
