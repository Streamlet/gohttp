package gohttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
)

type HttpServer interface {
	Serve() error
	Shutdown() error
}

func NewTcpServer(address string, handler http.Handler, errorChan chan error) HttpServer {
	server := new(tcpServer)
	server.Handler = handler
	server.address = address
	server.errorChan = errorChan
	return server
}

func NewUnixServer(socketFile string, handler http.Handler, errorChan chan error) HttpServer {
	server := new(unixServer)
	server.Handler = handler
	server.socketFile = socketFile
	server.errorChan = errorChan
	return server
}

type tcpServer struct {
	http.Server
	address   string
	errorChan chan error
}

type unixServer struct {
	http.Server
	socketFile string
	errorChan  chan error
}

func serve(s *http.Server, l net.Listener, errorChan chan error) {
	go func() {
		err := s.Serve(l)
		if !errors.Is(err, http.ErrServerClosed) {
			errorChan <- err
		}
	}()
}

func shutdown(s *http.Server) error {
	return s.Shutdown(context.Background())
}

func (s *tcpServer) Serve() error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	serve(&s.Server, l, s.errorChan)
	return nil
}

func (s *tcpServer) Shutdown() error {
	return shutdown(&s.Server)
}

func (s *unixServer) Serve() error {
	_ = os.Remove(s.socketFile)
	l, err := net.Listen("unix", s.socketFile)
	if err != nil {
		return err
	}
	if err := os.Chmod(s.socketFile, 0666); err != nil {
		return err
	}
	serve(&s.Server, l, s.errorChan)
	return nil
}

func (s *unixServer) Shutdown() error {
	err := shutdown(&s.Server)
	_ = os.Remove(s.socketFile)
	return err
}
