package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	HTTPServer *http.Server
}

func (s *Server) Run() error {
	return s.HTTPServer.ListenAndServe()
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		HTTPServer: &http.Server{
			Addr:           ":" + port,
			Handler:        handler,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
