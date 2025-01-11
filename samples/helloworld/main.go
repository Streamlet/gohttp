package main

import (
	"github.com/Streamlet/gohttp"
)

func HelloWorld(c gohttp.HttpContext) {
	c.String("Hello, World!")
}

func main() {
	application := gohttp.NewApplication[gohttp.HttpContext](gohttp.NewDefaultFactory(nil))
	application.Handle("/", HelloWorld)
	application.ServePort(80)
}
