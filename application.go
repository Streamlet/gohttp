package gohttp

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Application[T HttpContext] interface {
	ServeSock(sock string)
	ServePort(port uint)
	Handle(pattern string, handler func(T))
}

func NewApplication[T HttpContext](contextFactory ContextFactory[T], cacheProvider CacheProvider) Application[T] {
	return &application[T]{NewRouter[T](cacheProvider, contextFactory)}
}

type application[T HttpContext] struct {
	RouterInterface[T]
}

func (app *application[T]) ServeSock(sock string) {
	errorChan := make(chan error)
	srv := NewSockServer(sock, app, errorChan)
	app.serveUntilExit(srv, errorChan, sock, 0)
}

func (app *application[T]) ServePort(port uint) {
	errorChan := make(chan error)
	srv := NewPortServer(port, app, errorChan)
	app.serveUntilExit(srv, errorChan, "", port)
}

func (app *application[T]) serveUntilExit(srv HttpServer, errorChan chan error, sock string, port uint) {
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
