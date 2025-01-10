package server

import (
	"gohttp/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Application[T any] struct {
	router *web.Router[T]
}

func NewApplication[T any](cacheProvider web.CacheProvider, extContext T) *Application[T] {
	app := &Application[T]{web.NewRouter[T](cacheProvider, extContext)}
	return app
}

func (app *Application[T]) ServeSock(sock string) {
	errorChan := make(chan error)
	srv := NewSockServer(sock, app.router, errorChan)
	app.serveUntilExit(srv, errorChan, sock, 0)
}

func (app *Application[T]) ServePort(port uint) {
	errorChan := make(chan error)
	srv := NewPortServer(port, app.router, errorChan)
	app.serveUntilExit(srv, errorChan, "", port)
}

func (app *Application[T]) serveUntilExit(srv HttpServer, errorChan chan error, sock string, port uint) {
	err := srv.Serve()
	if err != nil {
		log.Printf("Failed to start server: %s.", err.Error())
		return
	}

	if sock != "" {
		log.Printf("Server started on sock '%s'.", sock)
	} else if port > 0 {
		log.Printf("Server started on port %d.", port)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			log.Printf("Server shutdown.\n")
		case err := <-errorChan:
			log.Printf("Server error: %s\n", err.Error())
		}
		break
	}

	_ = srv.Shutdown()
}

func (app *Application[T]) Handle(pattern string, handler func(web.Context[T])) {
	app.router.Handle(pattern, handler)
}
