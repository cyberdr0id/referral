// Package handler responsible for rounting.
package handler

import (
	"github.com/cyberdr0id/cv-web-service/internal/controllers"
	"github.com/gorilla/mux"
)

// InitRoutes initialize all endpoints.
func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", controllers.LogIn).Methods("POST")
	r.HandleFunc("/auth/signup", controllers.SignUp).Methods("POST")

	r.HandleFunc("/cv/send", controllers.SendCV).Methods("POST")
	r.HandleFunc("/cv/requests", controllers.GetRequests).Methods("GET")
	r.HandleFunc("/cv/load/{id}", controllers.LoadCV).Methods("GET")
	r.HandleFunc("/cv/filter/{type}", controllers.FilterCV).Methods("GET")

	return r
}
