package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/cyberdr0id/referral/internal/service"
	"github.com/gorilla/mux"
)

// Server presents a type of main application server.
type Server struct {
	HTTPServer *http.Server
	Router     *mux.Router
	Auth       service.Auth
	Referral   service.Referral
}

// Run starts server.
func (s *Server) Run(port string, handler http.Handler) error {
	s.HTTPServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.HTTPServer.ListenAndServe()
}

// NewServer creates a new instance of type Server.
func NewServer(auth service.Auth, referral service.Referral) *Server {
	s := &Server{
		Router:   mux.NewRouter(),
		Auth:     auth,
		Referral: referral,
	}

	s.InitRoutes()

	return s
}

// ServeHTTP dispatches the handler registered in the matched route.
func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(rw, r)
}

// Shutdown shuts down the server without interrupting any active connections.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
