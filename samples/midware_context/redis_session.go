package main

import (
	"context"
	"github.com/Streamlet/gohttp"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewSessionProvider(client *redis.Client, sessionKeyPrefix string) gohttp.CacheProvider {
	return &redisCache{client, sessionKeyPrefix}
}

type redisCache struct {
	client           *redis.Client
	sessionKeyPrefix string
}

func (rc *redisCache) Exists(key string) bool {
	if r, err := rc.client.Exists(context.Background(), rc.sessionKeyPrefix+key).Result(); err == nil && r > 0 {
		return true
	} else {
		return false
	}
}

func (rc *redisCache) HExists(key, field string) bool {
	if r, err := rc.client.HExists(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && r {
		return true
	} else {
		return false
	}
}

func (rc *redisCache) HGet(key, field string) interface{} {
	if r, err := rc.client.HGet(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil {
		return r
	} else {
		return nil
	}
}

func (rc *redisCache) HSet(key, field string, value interface{}, expiration time.Duration) {
	if r, err := rc.client.HSet(context.Background(), rc.sessionKeyPrefix+key, field, value).Result(); err != nil || r <= 0 {
		return
	}

	if expiration > 0 {
		if r, err := rc.client.Expire(context.Background(), rc.sessionKeyPrefix+key, expiration).Result(); err != nil || !r {
			return
		}
	}
}

func (rc *redisCache) HDelete(key, field string) bool {
	if r, err := rc.client.HExists(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && !r {
		return true
	}
	if r, err := rc.client.HDel(context.Background(), rc.sessionKeyPrefix+key, field).Result(); err == nil && r > 0 {
		return true
	}
	return false
}
