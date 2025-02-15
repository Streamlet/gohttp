package main

import (
	"github.com/Streamlet/gohttp"
)

func main() {
	application := gohttp.NewApplication[HttpContext](NewContextFactory())
	application.Handle("/", SuccessHandler)
	application.Handle("/success", SuccessWithDataHandler)
	application.Handle("/error", ErrorHandler)
	application.Handle("/error_message", ErrorWithMessageHandler)
	application.ServeTcp(":80")
}
