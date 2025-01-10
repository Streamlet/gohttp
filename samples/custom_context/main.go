package main

import (
	"gohttp/server"
	"gohttp/web"
)

type CustomContext interface {
	CustomFunc() string
}

type customContext struct {
}

func (cc customContext) CustomFunc() string {
	return "custom_func"
}

func CustomHandler(c web.Context[CustomContext]) {
	c.String(c.Custom().CustomFunc())
}

func main() {
	application := server.NewApplication[CustomContext](nil, &customContext{})
	application.Handle("/", CustomHandler)
	application.ServePort(80)
}
