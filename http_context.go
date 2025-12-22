package gohttp

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type HttpContext interface {
	io.Closer
	BasicContext
	RequestReader
	ResponseWriter
}

type BasicContext interface {
	HttpRequest() *http.Request
	HttpResponseWriter() http.ResponseWriter
	Session() Session
}

type Session interface {
	Id() string
	Exists(key string) bool
	Get(key string) interface{}
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string) bool
}

type RequestReader interface {
	GetQueryStrings() map[string][]string
	GetQueryStringValues(key string) []string
	GetQueryStringValue(key string) string
	GetRequestBodyAsBytes() ([]byte, error)
	GetRequestBodyAsString() (string, error)
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

func NewHttpContext(w http.ResponseWriter, r *http.Request, sm SessionManager, cookieDomain string) HttpContext {
	return &httpContext{
		responseWriter: w,
		request:        r,
		sessionManager: sm,
		cookieDomain:   cookieDomain,
	}
}

type httpContext struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	sessionManager SessionManager
	cookieDomain   string
}

func (c *httpContext) Close() error {
	return nil
}

func (c *httpContext) HttpRequest() *http.Request {
	return c.request
}

func (c *httpContext) HttpResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

const CookieSession = "SESSION"

func (c *httpContext) Session() Session {
	var session Session
	cookie, err := c.request.Cookie(CookieSession)
	if err == nil {
		session = c.sessionManager.GetSession(cookie.Value)
	}
	var sid string
	if session == nil {
		sid, session = c.sessionManager.CreateSession()
		cookie := http.Cookie{
			Name:     CookieSession,
			Value:    sid,
			Path:     "/",
			Domain:   c.cookieDomain,
			Secure:   c.request.TLS != nil || c.request.Header.Get("X-Forwarded-Proto") == "https",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		if c.cookieDomain != "" {
			if c.cookieDomain[0] == '~' {
				if re, err := regexp.Compile(c.cookieDomain[1:]); err == nil {
					cookie.Domain = re.FindString(c.request.Host)
				} else {
					log.Panicln("failed to compile regexp:", c.cookieDomain[1:]))
				}
			} else {
				cookie.Domain = c.cookieDomain
			}
		}
		http.SetCookie(c.responseWriter, &cookie)
	}
	return session
}

func (c *httpContext) GetQueryStrings() map[string][]string {
	return c.request.URL.Query()
}

func (c *httpContext) GetQueryStringValues(key string) []string {
	v, ok := c.GetQueryStrings()[key]
	if !ok {
		return nil
	}
	return v
}

func (c *httpContext) GetQueryStringValue(key string) string {
	v := c.GetQueryStringValues(key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (c *httpContext) GetRequestBodyAsBytes() ([]byte, error) {
	b, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *httpContext) GetRequestBodyAsString() (string, error) {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *httpContext) GetRequestBodyAsXml(v interface{}) error {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(bytes, v)
}

func (c *httpContext) GetRequestBodyAsJson(v interface{}) error {
	bytes, err := c.GetRequestBodyAsBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (c *httpContext) HttpError(statusCode int) {
	c.responseWriter.WriteHeader(statusCode)
}

func (c *httpContext) Redirect(url string) {
	c.responseWriter.Header().Add("Location", url)
	c.responseWriter.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *httpContext) String(response string) {
	c.responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	c.responseWriter.WriteHeader(http.StatusOK)
	_, err := c.responseWriter.Write([]byte(response))
	if err != nil {
		log.Println("failed to write response", err.Error(), response)
	}
}

func (c *httpContext) Xml(r interface{}) {
	xmlBytes, err := xml.Marshal(r)
	if err != nil {
		log.Println("failed to encode xml", err.Error(), r)
		c.responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.responseWriter.Header().Add("Content-Type", "application/xml; charset=utf-8")
	c.responseWriter.WriteHeader(http.StatusOK)
	_, err = c.responseWriter.Write([]byte(xml.Header))
	if err != nil {
		log.Println("failed to write response", err.Error())
	}
	_, err = c.responseWriter.Write(xmlBytes)
	if err != nil {
		log.Println("failed to write response", err.Error(), string(xmlBytes))
	}
}

func (c *httpContext) Json(r interface{}) {
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
