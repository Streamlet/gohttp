package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
)

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rc.Ping(context.Background()).Result()
	if err != nil {
		_ = rc.Close()
		log.Print("failed to connect to redis: ", err.Error())
		return
	}
	mysql, err := sql.Open("mysql", "root@tcp(localhost)/mysql?charset=latin1&loc=Local&parseTime=True&clientFoundRows=true")
	if err != nil {
		log.Print("failed to connect to mysql: ", err.Error())
		return
	}
	application := gohttp.NewApplication[HttpContext](NewContextFactory(rc, mysql))
	application.Handle("/redis_get", RedisGetHandler)
	application.Handle("/redis_set", RedisSetHandler)
	application.Handle("/mysql", DbHandler)
	application.ServePort(80)

	_ = mysql.Close()
	_ = rc.Close()
}
