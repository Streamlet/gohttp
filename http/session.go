package http

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

type Session interface {
	Exists(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string) bool
}

type CacheProvider interface {
	Exists(key string) bool
	HExists(key, field string) bool
	HGet(key, field string) interface{}
	HSet(key, field string, value interface{}, expiration time.Duration)
	HDelete(key, field string) bool
}

func GetSession(sessionId string, provider CacheProvider) Session {
	if provider.Exists(sessionId) {
		return newSession(sessionId, provider)
	} else {
		return nil
	}
}

func CreateSession(provider CacheProvider) (string, Session) {
	sessionId := generateSessionId()
	return sessionId, newSession(sessionId, provider)
}

func generateSessionId() string {
	b := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, b)
	return base64.URLEncoding.EncodeToString(b)
}

type session struct {
	sessionId string
	provider  CacheProvider
}

func newSession(sessionId string, provider CacheProvider) Session {
	return &session{sessionId, provider}
}

func (s *session) Exists(key string) bool {
	return s.provider.HExists(s.sessionId, key)
}

func (s *session) Get(key string) interface{} {
	return s.provider.HGet(s.sessionId, key)
}

func (s *session) Set(key string, value interface{}, expiration time.Duration) {
	s.provider.HSet(s.sessionId, key, value, expiration)
}

func (s *session) Delete(key string) bool {
	if !s.provider.HExists(s.sessionId, key) {
		return true
	}
	if s.provider.HDelete(s.sessionId, key) {
		return true
	}
	return false
}
