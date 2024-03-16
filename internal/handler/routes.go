package handler

import (
	"github.com/gorilla/mux"
	"vk/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() *mux.Router {
	r := mux.NewRouter()
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sing-up", h.SingUp).Methods("POST")
	auth.HandleFunc("/sing-in", h.SingIn).Methods("POST")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(h.UserIdentify)
	api.HandleFunc("/", Hello).Methods("GET")

	actors := api.PathPrefix("/actors").Subrouter()
	actors.HandleFunc("/", h.GetALLActors).Methods("GET")
	actors.HandleFunc("/{id}", h.GetActor).Methods("GET")
	actors.HandleFunc("/", h.SaveActor).Methods("POST")
	actors.HandleFunc("/{id}", UpdateActor).Methods("PATCH")
	actors.HandleFunc("/{id}", DeleteActor).Methods("DELETE")

	films := api.PathPrefix("/films").Subrouter()
	films.HandleFunc("/", h.GetALLFilms).Methods("GET")
	films.HandleFunc("/{id}", h.GetFilm).Methods("GET")
	films.HandleFunc("/", h.SaveFilm).Methods("POST")
	films.HandleFunc("/{id}", h.UpdateFilm).Methods("PATCH")
	films.HandleFunc("/{id}", h.DeleteFilm).Methods("DELETE")

	return r
}
