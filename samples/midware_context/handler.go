package main

import "context"

func RedisSetHandler(c HttpContext) {
	c.Cache().Set(context.Background(), "foo", "bar", 0)
}

func RedisGetHandler(c HttpContext) {
	r := c.Cache().Get(context.Background(), "foo")
	c.String(r.String())
}

func DbHandler(c HttpContext) {
	rs, err := c.DB().Query("select Db from mysql.db;")
	if err != nil {
		c.String(err.Error())
		return
	}

	dbNames := []string{}
	for rs.Next() {
		dbName := ""
		_ = rs.Scan(&dbName)
		dbNames = append(dbNames, dbName)
	}
	c.Json(dbNames)
}
