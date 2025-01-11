package main

func CustomContextHandler(c HttpContext) {
	c.String(c.Custom())
}
