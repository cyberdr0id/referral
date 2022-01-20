// Package handler responsible for rounting.
package handler

// InitRoutes initialize all endpoints.
func (s *Server) InitRoutes() {
	s.Router.HandleFunc("/auth/login", s.LogIn).Methods("POST")
	s.Router.HandleFunc("/auth/signup", s.SignUp).Methods("POST")

	userRouter := s.Router.NewRoute().Subrouter()
	userRouter.Use(s.AuthorizationMiddleware)

	userRouter.HandleFunc("/references", s.SendCandidate).Methods("POST")
	userRouter.HandleFunc("/references", s.GetRequests).Methods("GET")
	userRouter.HandleFunc("/cvs", s.DownloadCV).Methods("GET")

	adminRouter := userRouter.NewRoute().Subrouter()
	adminRouter.Use(s.AdminMiddleware)

	adminRouter.HandleFunc("/admin/references", s.UpdateRequest).Methods("PUT")
	adminRouter.HandleFunc("/admin/references", s.GetAllRequests).Methods("GET")
}
