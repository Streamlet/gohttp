package main

func SuccessHandler(c HttpContext) {
	c.Success(nil)
}

func SuccessWithDataHandler(c HttpContext) {
	type SuccessData struct {
		Value   int    `json:"value"`
		Message string `json:"message"`
	}
	c.Success(SuccessData{123, "ok"})
}

func ErrorHandler(c HttpContext) {
	c.Error(ErrorInternal, "")
}

func ErrorWithMessageHandler(c HttpContext) {
	c.Error(ErrorInternal, "internal error")
}
