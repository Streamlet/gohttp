package main

import (
	"net/http"

	"github.com/Streamlet/gohttp"
)

type HttpContext interface {
	gohttp.HttpContext
	Success(data interface{})
	Error(errorCode int, errorMessage string)
}

func NewContextFactory() gohttp.ContextFactory[HttpContext] {
	return &contextFactory{gohttp.NewSessionManager(nil)}
}

type contextFactory struct {
	sm gohttp.SessionManager
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &httpContext{
		gohttp.NewHttpContext(w, r, cf.sm, ""),
	}
}

type httpContext struct {
	gohttp.HttpContext
}

const (
	ErrorSuccess  int = 0
	ErrorInternal int = -1
)

type response struct {
	Success bool        `json:"success,omitempty"`
	Error   int         `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (c *httpContext) Success(data interface{}) {
	c.Json(response{true, ErrorSuccess, data, ""})
}

func (c *httpContext) Error(errorCode int, errorMessage string) {
	c.Json(response{false, errorCode, nil, errorMessage})
}
