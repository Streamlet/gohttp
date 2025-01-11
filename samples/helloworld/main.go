package main

import (
	"github.com/Streamlet/gohttp"
)

func main() {
	application := gohttp.NewApplication[gohttp.HttpContext](gohttp.NewDefaultFactory(nil))
	application.Handle("/", HelloWorld)
	application.ServePort(80)
}
