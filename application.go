package gohttp

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Application[T HttpContext] interface {
	ServeUnix(socketFile string)
	ServeTcp(address string)
	RawHandle(pattern string, handler http.Handler)
	Handle(pattern string, handler func(T))
}

func NewApplication[T HttpContext](contextFactory ContextFactory[T]) Application[T] {
	return &application[T]{NewRouter[T](contextFactory)}
}

type application[T HttpContext] struct {
	RouterInterface[T]
}

func (app *application[T]) ServeUnix(socketFile string) {
	errorChan := make(chan error)
	srv := NewUnixServer(socketFile, app, errorChan)
	app.serveUntilExit(srv, errorChan, socketFile, "")
}

func (app *application[T]) ServeTcp(address string) {
	errorChan := make(chan error)
	srv := NewTcpServer(address, app, errorChan)
	app.serveUntilExit(srv, errorChan, "", address)
}

func (app *application[T]) serveUntilExit(srv HttpServer, errorChan chan error, socketFile string, address string) {
	err := srv.Serve()
	if err != nil {
		log.Printf("Failed to start server: %s.", err.Error())
		return
	}

	if socketFile != "" {
		log.Printf("Server started on socket file '%s'.", socketFile)
	} else if address != "" {
		log.Printf("Server started on address '%s'.", address)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		log.Printf("Server shutdown.\n")
	case err := <-errorChan:
		log.Printf("Server error: %s\n", err.Error())
	}

	_ = srv.Shutdown()
}
