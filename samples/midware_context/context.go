package main

import (
	"database/sql"
	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type HttpContext interface {
	gohttp.HttpContext
	Cache() *redis.Client
	DB() *sql.DB
}

func NewContextFactory(cache *redis.Client, db *sql.DB) gohttp.ContextFactory[HttpContext] {
	return &contextFactory{gohttp.NewSessionManager(NewSessionProvider(cache, "SESSION_")), cache, db}
}

type contextFactory struct {
	sm    gohttp.SessionManager
	cache *redis.Client
	db    *sql.DB
}

func (cf *contextFactory) NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &httpContext{
		gohttp.NewHttpContext(w, r, cf.sm), cf.cache, cf.db,
	}
}

type httpContext struct {
	gohttp.HttpContext
	cache *redis.Client
	db    *sql.DB
}

func (c *httpContext) Cache() *redis.Client {
	return c.cache
}

func (c *httpContext) DB() *sql.DB {
	return c.db
}
