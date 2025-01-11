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

func (s *redisCache) Exists(key string) bool {
	if r, err := s.client.Exists(context.Background(), s.sessionKeyPrefix+key).Result(); err == nil && r > 0 {
		return true
	} else {
		return false
	}
}

func (s *redisCache) HExists(key, field string) bool {
	if r, err := s.client.HExists(context.Background(), s.sessionKeyPrefix+key, field).Result(); err == nil && r {
		return true
	} else {
		return false
	}
}

func (s *redisCache) HGet(key, field string) interface{} {
	if r, err := s.client.HGet(context.Background(), s.sessionKeyPrefix+key, field).Result(); err == nil {
		return r
	} else {
		return nil
	}
}

func (s *redisCache) HSet(key, field string, value interface{}, expiration time.Duration) {
	if r, err := s.client.HSet(context.Background(), s.sessionKeyPrefix+key, field, value).Result(); err != nil || r <= 0 {
		return
	}

	if expiration > 0 {
		if r, err := s.client.Expire(context.Background(), s.sessionKeyPrefix+key, expiration).Result(); err != nil || !r {
			return
		}
	}
}

func (s *redisCache) HDelete(key, field string) bool {
	if r, err := s.client.HExists(context.Background(), s.sessionKeyPrefix+key, field).Result(); err == nil && !r {
		return true
	}
	if r, err := s.client.HDel(context.Background(), s.sessionKeyPrefix+key, field).Result(); err == nil && r > 0 {
		return true
	}
	return false
}
