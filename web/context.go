package web

import (
	"encoding/json"
	"encoding/xml"
	"gohttp/cache"
	"io"
	"log"
	"net/http"
)

type RequestReader interface {
	GetQueryStrings() map[string][]string
	GetQueryStringValues(key string) []string
	GetQueryStringValue(key string) string
	GetRequestBodyAsBytes() ([]byte, error)
	GetRequestBodyAsStrings() (string, error)
	GetRequestBodyAsXml(v interface{}) error
	GetRequestBodyAsJson(v interface{}) error
}

type ResponseWriter interface {
	HttpError(statusCode int)
	Redirect(url string)
	String(response string)
	Xml(r interface{})
	Json(r interface{})
}

type BasicContext interface {
	HttpRequest() *http.Request
	HttpResponseWriter() http.ResponseWriter
	Session() Session
}

type Context[T any] interface {
	BasicContext
	RequestReader
	ResponseWriter
	Custom() T
}

type httpContext[T any] struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	cacheProvider  CacheProvider
	session        Session
	customContext  T
}

func newHttpContext[T any](w http.ResponseWriter, r *http.Request, cp CacheProvider, cc T) *httpContext[T] {
	if cp == nil {
		cp = cache.NewMemoryCache()
	}
	return &httpContext[T]{
		responseWriter: w,
		request:        r,
		cacheProvider:  cp,
		customContext:  cc,
	}
}

func (c *httpContext[T]) HttpRequest() *http.Request {
	return c.request
}

func (c *httpContext[T]) HttpResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

const CookieSession = "SESSION"

func (c *httpContext[T]) Session() Session {
	if c.session == nil && c.cacheProvider != nil {
		var s Session
		cookie, err := c.request.Cookie(CookieSession)
		if err == nil {
			s = GetSession(cookie.Value, c.cacheProvider)
		}
		var sid string
		if s == nil {
			sid, s = CreateSession(c.cacheProvider)
			cookie := http.Cookie{Name: CookieSession, Value: sid, Path: "/"}
			if c.request.URL.Scheme != "https" {
				cookie.SameSite = http.SameSiteNoneMode
				cookie.Secure = true
			}
			http.SetCookie(c.responseWriter, &cookie)
		}
		c.session = s
	}
	return c.session
}

func (c *httpContext[T]) Release() {
	if c.session != nil {
		c.session = nil
	}
}

func (c *httpContext[T]) Custom() T {
	return c.customContext
}

func (c *httpContext[T]) GetQueryStrings() map[string][]string {
	return c.request.URL.Query()
}

func (c *httpContext[T]) GetQueryStringValues(key string) []string {
	v, ok := c.GetQueryStrings()[key]
	if !ok {
		return nil
	}
	return v
}

func (c *httpContext[T]) GetQueryStringValue(key string) string {
	v := c.GetQueryStringValues(key)
	if v == nil || len(v) == 0 {
		return ""
	}
	return v[0]
}

func (c *httpContext[T]) GetRequestBodyAsBytes() ([]byte, error) {
	b, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *httpContext[T]) GetRequestBodyAsStrings() (string, error) {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *httpContext[T]) GetRequestBodyAsXml(v interface{}) error {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(bytes, v)
}

func (c *httpContext[T]) GetRequestBodyAsJson(v interface{}) error {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (c *httpContext[T]) HttpError(statusCode int) {
	c.responseWriter.WriteHeader(statusCode)
}

func (c *httpContext[T]) Redirect(url string) {
	c.responseWriter.Header().Add("Location", url)
	c.responseWriter.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *httpContext[T]) String(response string) {
	c.responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	c.responseWriter.WriteHeader(http.StatusOK)
	_, err := c.responseWriter.Write([]byte(response))
	if err != nil {
		log.Println("failed to write response", err.Error(), response)
	}
}

func (c *httpContext[T]) Xml(r interface{}) {
	xmlString, err := xml.Marshal(r)
	if err != nil {
		log.Println("failed to encode xml", err.Error(), r)
		c.responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.responseWriter.Header().Add("Content-Type", "application/xml; charset=utf-8")
	c.responseWriter.WriteHeader(http.StatusOK)
	_, err = c.responseWriter.Write(xmlString)
	if err != nil {
		log.Println("failed to write response", err.Error(), string(xmlString))
	}
}

func (c *httpContext[T]) Json(r interface{}) {
	jsonString, err := json.Marshal(r)
	if err != nil {
		log.Println("failed to encode json", err.Error(), r)
		c.responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.responseWriter.Header().Add("Content-Type", "application/json; charset=utf-8")
	c.responseWriter.WriteHeader(http.StatusOK)
	_, err = c.responseWriter.Write(jsonString)
	if err != nil {
		log.Println("failed to write response", err.Error(), jsonString)
	}
}
