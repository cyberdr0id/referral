// Package handler responsible for rounting.
package handler

// InitRoutes initialize all endpoints.
func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/auth/login", s.LogIn).Methods("POST")   // user authorization
	s.Router.HandleFunc("/auth/signup", s.SignUp).Methods("POST") // user registration

	s.Router.HandleFunc("/references", s.SendCandidate).Methods("POST") // sending candidate
	s.Router.HandleFunc("/references", s.GetRequests).Methods("GET")    // user request history
	s.Router.HandleFunc("/cvs", s.DownloadCV).Methods("GET")            // loading cv
}
