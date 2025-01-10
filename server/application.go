package server

import (
	"gohttp/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Application struct {
	router *http.Router
}

func NewApplication(cacheProvider http.CacheProvider) *Application {
	app := &Application{http.NewRouter(cacheProvider)}
	return app
}

func (app *Application) ServeSock(sock string) {
	errorChan := make(chan error)
	srv := NewSockServer(sock, app.router, errorChan)
	app.serveUntilExit(srv, errorChan, sock, 0)
}

func (app *Application) ServePort(port uint) {
	errorChan := make(chan error)
	srv := NewPortServer(port, app.router, errorChan)
	app.serveUntilExit(srv, errorChan, "", port)
}

func (app *Application) serveUntilExit(srv HttpServer, errorChan chan error, sock string, port uint) {
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

func (app *Application) Handle(pattern string, handler func(http.Context)) {
	app.router.Handle(pattern, handler)
}
