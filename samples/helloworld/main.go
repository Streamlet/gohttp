package main

import (
	"github.com/Streamlet/gohttp"
)

func main() {
	application := gohttp.NewApplication[gohttp.HttpContext](gohttp.NewContextFactory(nil, ""))
	application.Handle("/", HelloWorld)
	application.ServeTcp(":80")
}
