// Package handler responsible for rounting.
package handler

import (
	"github.com/gorilla/mux"
)

// InitRoutes initialize all endpoints.
func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", LogIn).Methods("POST")   // user authorization
	r.HandleFunc("/auth/signup", SignUp).Methods("POST") // user registration

	r.HandleFunc("/references", SendCandidate).Methods("POST") // sending candidate
	r.HandleFunc("/references", GetRequests).Methods("GET")    // user request history
	r.HandleFunc("/cvs/{id}", DownloadCV).Methods("GET")       // loading cv

	return r
}
