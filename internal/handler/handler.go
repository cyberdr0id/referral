package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("POST")
	r.HandleFunc("/auth/signup", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("POST")

	r.HandleFunc("/cv/send", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("POST")
	r.HandleFunc("/cv/requests", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("GET")
	r.HandleFunc("/cv/load/{id}", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("GET")
	r.HandleFunc("/cv/filter/{type}", func(rw http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	return r
}
