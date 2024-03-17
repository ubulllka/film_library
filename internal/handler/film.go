package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"vk/internal/models/DTO"
)

func (h *Handler) GetALLFilms(w http.ResponseWriter, r *http.Request) {
	column := r.URL.Query().Get("column")
	order := r.URL.Query().Get("order")

	films, err := h.service.Film.GetAllFilms(column, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filmsByte, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsByte)
}

func (h *Handler) GetFilm(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["id"]

	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	film, err := h.service.Film.GetFilm(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filmByte, err := json.Marshal(film)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmByte)
}

func (h *Handler) SearchFilms(w http.ResponseWriter, r *http.Request) {
	fragment := r.URL.Query().Get("fragment")

	films, err := h.service.Film.SearchFilms(fragment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filmsByte, err := json.Marshal(films)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsByte)
}

func (h *Handler) SaveFilm(w http.ResponseWriter, r *http.Request) {
	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}

	var filmSave DTO.FilmInput
	if err := json.NewDecoder(r.Body).Decode(&filmSave); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	film, actors, err := filmSave.GetFilmAndActors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(film); err != nil {
		errs := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errs), http.StatusBadRequest)
		return
	}

	id, err := h.service.Film.CreateFilm(film, actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Film %d create\"}", id))))
}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}

	var input DTO.FilmUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Film.UpdateFilm(id, input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Film %d update\"}", id))))
}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}

	if err := h.service.Film.DeleteFilm(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Delete %d delete\"}", id))))
}
