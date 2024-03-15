package handler

import (
	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", Hello).Methods("GET")
	return r
}
