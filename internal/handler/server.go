package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/cyberdr0id/referral/internal/service"
	"github.com/gorilla/mux"
)

type Server struct {
	HTTPServer *http.Server
	Router     *mux.Router
	Auth       service.Auth
	Referral   service.Referral
}

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

func NewServer(auth service.Auth, referral service.Referral) *Server {
	s := &Server{
		Router:   mux.NewRouter(),
		Auth:     auth,
		Referral: referral,
	}

	s.InitRoutes()

	return s
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(rw, r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
