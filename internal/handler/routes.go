package handler

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger/v2"
	_ "vk/docs"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sing-up", h.SingUp).Methods("POST")
	auth.HandleFunc("/sing-in", h.SingIn).Methods("POST")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(h.UserIdentify)

	actors := api.PathPrefix("/actors").Subrouter()
	actors.HandleFunc("", h.GetALLActors).Methods("GET")
	actors.HandleFunc("/{id}", h.GetActor).Methods("GET")
	actors.HandleFunc("", h.SaveActor).Methods("POST")
	actors.HandleFunc("/{id}", h.UpdateActor).Methods("PATCH")
	actors.HandleFunc("/{id}", h.DeleteActor).Methods("DELETE")

	films := api.PathPrefix("/films").Subrouter()
	films.HandleFunc("", h.GetALLFilms).Methods("GET")
	films.HandleFunc("/search", h.SearchFilms).Methods("GET")
	films.HandleFunc("/{id}", h.GetFilm).Methods("GET")
	films.HandleFunc("", h.SaveFilm).Methods("POST")
	films.HandleFunc("/{id}", h.UpdateFilm).Methods("PATCH")
	films.HandleFunc("/{id}", h.DeleteFilm).Methods("DELETE")

	return r
}
