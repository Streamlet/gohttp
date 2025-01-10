package main

import (
	"gohttp/server"
	"gohttp/web"
)

func HelloWorld(c web.Context) {
	c.String("Hello, World!")
}

func main() {
	application := server.NewApplication(nil, nil)
	application.Handle("/", HelloWorld)
	application.ServePort(80)
}
