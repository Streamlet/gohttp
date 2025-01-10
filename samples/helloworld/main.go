package main

import (
	"gohttp/server"
	"gohttp/web"
)

func HelloWorld(c web.Context[interface{}]) {
	c.String("Hello, World!")
}

func main() {
	application := server.NewApplication[interface{}](nil, nil)
	application.Handle("/", HelloWorld)
	application.ServePort(80)
}
