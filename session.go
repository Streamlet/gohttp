package gohttp

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"
)

const (
	DefaultSessionExpireMinutes = 30
)

type CacheProvider interface {
	Exists(key string) bool
	Touch(key string, expiration time.Duration)
	HExists(key, field string) bool
	HGet(key, field string) interface{}
	HSet(key, field string, value interface{}, expiration time.Duration)
	HDelete(key, field string) bool
}

type SessionManager interface {
	GetSession(sessionId string) Session
	CreateSession() (string, Session)
}

func NewSessionManager(cacheProvider CacheProvider) SessionManager {
	return &sessionManager{
		cacheProvider,
	}
}

type sessionManager struct {
	cacheProvider CacheProvider
}

func (sm *sessionManager) GetSession(sessionId string) Session {
	if sm.cacheProvider.Exists(sessionId) {
		sm.cacheProvider.Touch(sessionId, time.Minute*DefaultSessionExpireMinutes)
		return newSession(sessionId, sm.cacheProvider)
	} else {
		return nil
	}
}

func (sm *sessionManager) CreateSession() (string, Session) {
	sessionId := generateSessionId()
	sm.cacheProvider.Touch(sessionId, time.Minute*DefaultSessionExpireMinutes)
	return sessionId, newSession(sessionId, sm.cacheProvider)
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

func (s *session) Id() string {
	return s.sessionId
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
