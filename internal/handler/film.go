package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"vk/internal/models/DTO"
)

// GetAllFilms
// @Summary		Get All Films
// @Tags		films
// @Description	get all films
// @Accept		json
// @Produce		json
// @Param		column	query		string	false	"sort column"	Enums(film_name, release_date, rating)
// @Param		order	query		string	false	"sort order"	Enums(ASC, DESC)
// @Success		200		{array}		DTO.FilmDTO
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/films [get]
func (h *Handler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	column := r.URL.Query().Get("column")
	order := r.URL.Query().Get("order")

	films, err := h.service.Film.GetAllFilms(column, order)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	filmsByte, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("GetAllFilms with column - " + column + "; order - " + order)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsByte)
}

// SearchFilms
// @Summary		Search Films
// @Tags		films
// @Description	search films
// @Accept		json
// @Produce		json
// @Param		q		query		string	false	"search by q"
// @Success		200		{array}		DTO.FilmDTO
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/films/search [get]
func (h *Handler) SearchFilms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fragment := r.URL.Query().Get("q")
	films, err := h.service.Film.SearchFilms(fragment)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	filmsByte, err := json.Marshal(films)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("SearchFilms with q - " + fragment)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmsByte)
}

// GetFilm
// @Summary		Get Film
// @Tags		films
// @Description	get film by id
// @Accept		json
// @Produce		json
// @Param		id		path		int	true	"film id"
// @Success		200		{array}		DTO.FilmDTO
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/films/{id} [get]
func (h *Handler) GetFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]

	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	film, err := h.service.Film.GetFilm(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	filmByte, err := json.Marshal(film)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("GetFilm with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filmByte)
}

// SaveFilm
// @Summary		Save Film
// @Security	ApiKeyAuth
// @Tags		films
// @Description	create film
// @Accept		json
// @Produce		json
// @Param		filmSave	body		DTO.FilmInput	true	"film info"
// @Success		200			{integer}	integer			1
// @Failure		400,403		{object}	errorResponse
// @Failure		500			{object}	errorResponse
// @Failure		default		{object}	errorResponse
// @Router		/api/v1/films [post]
func (h *Handler) SaveFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	var filmSave DTO.FilmInput
	if err := json.NewDecoder(r.Body).Decode(&filmSave); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	film, actors, err := filmSave.GetFilmAndActors()
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(film); err != nil {
		errs := err.(validator.ValidationErrors)
		newErrorResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	id, err := h.service.Film.CreateFilm(film, actors)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("SaveFilm with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

// UpdateFilm
// @Summary		Update Film
// @Security	ApiKeyAuth
// @Tags		films
// @Description	update film
// @Accept		json
// @Produce		json
// @Param		id		path		int				true	"film id"
// @Param		input	body		DTO.FilmUpdate	true	"film update info"
// @Success		200		{object}	statusResponse
// @Failure		400,403	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/films/{id} [patch]
func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	var input DTO.FilmUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Film.UpdateFilm(id, input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("UpdateFilm with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(statusResponse{"ok"})
	_, _ = w.Write(result)
}

// DeleteFilm
// @Summary		Delete Film
// @Security	ApiKeyAuth
// @Tags		films
// @Description	delete film
// @Accept		json
// @Produce		json
// @Param		id		path		int	true	"film id"
// @Success		200		{object}	statusResponse
// @Failure		400,403	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/films/{id} [delete]
func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userRole := context.Get(r, USERROLE)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	if err := h.service.Film.DeleteFilm(id); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("DeleteFilm with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(statusResponse{"ok"})
	_, _ = w.Write(result)
}
