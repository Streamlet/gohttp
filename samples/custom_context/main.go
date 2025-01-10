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

func CustomHandler(c web.Context) {
	c.String(c.Custom().(CustomContext).CustomFunc())
}

func main() {
	application := server.NewApplication(nil, &customContext{})
	application.Handle("/", CustomHandler)
	application.ServePort(80)
}
