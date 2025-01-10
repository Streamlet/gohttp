package cache

import "time"

type MemoryCache interface {
	Exists(key string) bool
	HExists(key, field string) bool
	HGet(key, field string) interface{}
	HSet(key, field string, value interface{}, expiration time.Duration)
	HDelete(key, field string) bool
}

type cacheItem struct {
	data       interface{}
	createTime time.Time
	expiration time.Duration
}

type memoryCache struct {
	cache map[string]map[string]cacheItem
}

func NewMemoryCache() MemoryCache {
	return &memoryCache{map[string]map[string]cacheItem{}}
}

func (s *memoryCache) Exists(key string) bool {
	_, ok := s.cache[key]
	return ok
}

func (s *memoryCache) HExists(key, field string) bool {
	session, ok := s.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	return ok
}

func (s *memoryCache) HGet(key, field string) interface{} {
	session, ok := s.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	return item.data
}

func (s *memoryCache) HSet(key, field string, value interface{}, expiration time.Duration) {
	s.cache[key][field] = cacheItem{value, time.Now(), expiration}
}

func (s *memoryCache) HDelete(key, field string) bool {
	session, ok := s.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	delete(session, field)
	return true
}
