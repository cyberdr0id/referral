// Package handler responsible for rounting.
package handler

// InitRoutes initialize all endpoints.
func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/auth/login", s.LogIn).Methods("POST")
	s.Router.HandleFunc("/auth/signup", s.SignUp).Methods("POST")

	sr := s.Router.NewRoute().Subrouter()
	sr.Use(s.AuthorizationMiddleware)

	sr.HandleFunc("/references", s.SendCandidate).Methods("POST")
	sr.HandleFunc("/references", s.GetRequests).Methods("GET")
	sr.HandleFunc("/cvs", s.DownloadCV).Methods("GET")

	adminRouter := sr.NewRoute().Subrouter()
	adminRouter.Use(s.AdminMiddleware)

	adminRouter.HandleFunc("/admin/references", s.UpdateRequest).Methods("PUT")
	adminRouter.HandleFunc("/admin/references", s.GetAllRequests).Methods("GET")
}
