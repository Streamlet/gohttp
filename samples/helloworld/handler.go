package main

import (
	"github.com/Streamlet/gohttp"
)

func HelloWorld(c gohttp.HttpContext) {
	c.String("Hello, World!")
}
