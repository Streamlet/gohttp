package main

import (
	"github.com/Streamlet/gohttp/server"
	"github.com/Streamlet/gohttp/web"
)

type ExtContext interface {
	ExtFunc() string
}

type extContext struct {
}

func (ex extContext) ExtFunc() string {
	return "ext_func"
}

func ExtContextHandler(c web.Context[ExtContext]) {
	c.String(c.Ext().ExtFunc())
}

func main() {
	application := server.NewApplication[ExtContext](nil, &extContext{})
	application.Handle("/", ExtContextHandler)
	application.ServePort(80)
}
