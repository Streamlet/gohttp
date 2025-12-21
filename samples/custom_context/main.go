package main

import (
	"github.com/Streamlet/gohttp"
)

func main() {
	application := gohttp.NewApplication[HttpContext](NewContextFactory())
	application.Handle("/", CustomContextHandler)
	application.ServeTcp(":80")
}
