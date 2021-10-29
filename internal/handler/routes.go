// Package handler responsible for rounting.
package handler

import (
	"github.com/cyberdr0id/referral/internal/controllers"
	"github.com/gorilla/mux"
)

// InitRoutes initialize all endpoints.
func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", controllers.LogIn).Methods("POST")   // user authorization
	r.HandleFunc("/auth/signup", controllers.SignUp).Methods("POST") // user registration

	r.HandleFunc("/references", controllers.SendCV).Methods("POST")     // sending cv
	r.HandleFunc("/references", controllers.GetRequests).Methods("GET") // user request history
	r.HandleFunc("/cvs/{id}", controllers.LoadCV).Methods("GET")        // loading cv

	return r
}
