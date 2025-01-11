package main

import (
	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type HttpContext interface {
	gohttp.HttpContext
	Cache() *redis.Client
	DB() *Connection
}

func NewContextFactory(cache *redis.Client, db *Connection) gohttp.ContextFactory[HttpContext] {
	return &contextFactory{gohttp.NewSessionManager(NewSessionProvider(cache, "SESSION_")), cache, db}
}

type contextFactory struct {
	sm    gohttp.SessionManager
	cache *redis.Client
	db    *Connection
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &httpContext{
		gohttp.NewHttpContext(w, r, cf.sm), cf.cache, cf.db,
	}
}

type httpContext struct {
	gohttp.HttpContext
	cache *redis.Client
	db    *Connection
}

func (c *httpContext) Cache() *redis.Client {
	return c.cache
}

func (c *httpContext) DB() *Connection {
	return c.db
}
