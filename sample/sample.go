package main

import (
	"gohttp/http"
	"gohttp/server"
)

func HelloWorld(c http.Context) {
	c.String("Hello, World!")
}

func main() {
	application := server.NewApplication(nil)
	application.Handle("/", HelloWorld)
	application.ServePort(80)
}
