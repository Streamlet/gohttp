package gohttp

import "time"

func newMemoryCache() CacheProvider {
	return &memoryCache{map[string]map[string]cacheItem{}}
}

type memoryCache struct {
	cache map[string]map[string]cacheItem
}

type cacheItem struct {
	data       interface{}
	createTime time.Time
	expiration time.Duration
}

func (mc *memoryCache) Exists(key string) bool {
	_, ok := mc.cache[key]
	return ok
}

func (mc *memoryCache) HExists(key, field string) bool {
	session, ok := mc.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if item.expiration > 0 && !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	return ok
}

func (mc *memoryCache) HGet(key, field string) interface{} {
	session, ok := mc.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if item.expiration > 0 && !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	return item.data
}

func (mc *memoryCache) HSet(key, field string, value interface{}, expiration time.Duration) {
	mc.cache[key][field] = cacheItem{value, time.Now(), expiration}
}

func (mc *memoryCache) HDelete(key, field string) bool {
	session, ok := mc.cache[key]
	if !ok {
		return false
	}
	item, ok := session[field]
	if !ok {
		return false
	}
	if item.expiration > 0 && !item.createTime.Add(item.expiration).After(time.Now()) {
		delete(session, field)
		return false
	}
	delete(session, field)
	return true
}
