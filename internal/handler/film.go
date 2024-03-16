package handler

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (h *Handler) GetALLFilms(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetFilm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = id
}

func (h *Handler) SaveFilm(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}
}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = id
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = id
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}
}
