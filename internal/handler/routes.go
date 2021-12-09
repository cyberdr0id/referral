// Package handler responsible for rounting.
package handler

// InitRoutes initialize all endpoints.
func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/auth/login", s.LogIn).Methods("POST")   // user authorization
	s.Router.HandleFunc("/auth/signup", s.SignUp).Methods("POST") // user registration

	sr := s.Router.NewRoute().Subrouter()
	sr.Use(s.AuthorizationMiddleware)

	sr.HandleFunc("/references", s.SendCandidate).Methods("POST") // sending candidate
	sr.HandleFunc("/references", s.GetRequests).Methods("GET")    // user request history
	sr.HandleFunc("/references", s.UpdateRequest).Methods("PUT")  // user request history
	sr.HandleFunc("/cvs", s.DownloadCV).Methods("GET")            // loading cv
}
