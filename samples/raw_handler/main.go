package main

import (
	"embed"
	_ "embed"
	"net/http"

	"github.com/Streamlet/gohttp"
)

//go:embed *.html
var static embed.FS

func main() {
	application := gohttp.NewApplication[gohttp.HttpContext](gohttp.NewContextFactory(nil, ""))
	application.RawHandle("/", http.FileServer(http.FS(static)))
	application.ServeTcp(":80")
}
