# gohttp

A simple HTTP server framework for golang.

## Getting start

See `samples/helloworld`

```go
package main

import (
  "github.com/Streamlet/gohttp"
)

func HelloWorld(c gohttp.HttpContext) {
  c.String("Hello, World!")
}

func main() {
  application := gohttp.NewApplication[gohttp.HttpContext](gohttp.NewContextFactory(nil))
  application.Handle("/", HelloWorld)
  application.ServePort(80)
}
```

## HttpContext

In gohttp, HTTP handler signature is `func (c gohttp.HttpContext)`.

HttpContext provides:

* Raw HTTP request and response
  * HttpRequest() *http.Request
  * HttpResponseWriter() http.ResponseWriter
* Session
  * Default to a memory cached session
  * Can be replaced to redis based session simply:
    * Implement a CacheProvider
    * Pass CacheProvider to NewContextFactory
* Input
  * GetQueryStrings() map[string][]string
  * GetQueryStringValues(key string) []string
  * GetQueryStringValue(key string) string
  * GetRequestBodyAsBytes() ([]byte, error)
  * GetRequestBodyAsStrings() (string, error)
  * GetRequestBodyAsXml(v interface{}) error
  * GetRequestBodyAsJson(v interface{}) error
* Output
  * HttpError(statusCode int)
  * Redirect(url string)
  * String(response string)
  * Xml(r interface{})
  * Json(r interface{})

# Extend HttpContext

See `samples/custom_context`, `samples/std_json_context` and `samples/midware_context`.

General steps:

1. Define a new HttpContext, with gohttp.HttpContext embedded, and other functions appended.
2. Implement self-defined HttpContext
3. Define a new ContextFactory to create self-defined HttpContext
4. Pass the new ContextFactory to gohttp.NewApplication
5. Use self-defined HttpContext for all handlers.

Typically, more helper functions for input and output can be append to HttpContext.
Middleware instances, e.g. redis, mysql, are expected to be added to HttpContext, too.
