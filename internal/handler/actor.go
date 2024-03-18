package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"vk/internal/models/DTO"
)

// GetAllActors
// @Summary		Get All Actors
// @Tags		actors
// @Description	get all actors
// @Accept		json
// @Produce		json
// @Success		200		{array}		DTO.ActorDTO
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/actors [get]
func (h *Handler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	actors, err := h.service.Actor.GetAllActors()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	actorsByte, err := json.Marshal(actors)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("GetAllActors is ok")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorsByte)
}

// GetActor
// @Summary		Get Actor
// @Tags		actors
// @Description	get actor by id
// @Accept		json
// @Produce		json
// @Param		id		path		int	true	"actor id"
// @Success		200		{array}		DTO.ActorDTO
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/actors/{id} [get]
func (h *Handler) GetActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]

	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	actor, err := h.service.Actor.GetActor(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	actorByte, err := json.Marshal(actor)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("GetActor with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorByte)
}

// SaveActor
// @Summary		Save Actor
// @Security	ApiKeyAuth
// @Tags		actors
// @Description	create actor
// @Accept		json
// @Produce		json
// @Param		actorSave	body		DTO.ActorInput	true	"actor info"
// @Success		200			{integer}	integer			1
// @Failure		400,403		{object}	errorResponse
// @Failure		500			{object}	errorResponse
// @Failure		default		{object}	errorResponse
// @Router		/api/v1/actors [post]
func (h *Handler) SaveActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//userRole := context.Get(r, USERROLE)
	userRole := r.Context().Value(USERROLE)
	fmt.Println(userRole)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	var actorSave DTO.ActorInput
	if err := json.NewDecoder(r.Body).Decode(&actorSave); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	actor, err := actorSave.GetActor()
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(actor); err != nil {
		errs := err.(validator.ValidationErrors)
		newErrorResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	id, err := h.service.Actor.CreateActor(actor)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("SaveActor with id " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))

}

// UpdateActor
// @Summary		Update Actor
// @Security	ApiKeyAuth
// @Tags		actors
// @Description	update actor
// @Accept		json
// @Produce		json
// @Param		id		path		int				true	"actor id"
// @Param		input	body		DTO.ActorUpdate	true	"actor update info"
// @Success		200		{object}	statusResponse
// @Failure		400,403	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/actors/{id} [patch]
func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//userRole := context.Get(r, USERROLE)
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	var input DTO.ActorUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Actor.UpdateActor(id, input); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("UpdateActor with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(statusResponse{"ok"})
	_, _ = w.Write(result)
}

// DeleteActor
// @Summary		Delete Actor
// @Security	ApiKeyAuth
// @Tags		actors
// @Description	delete actor
// @Accept		json
// @Produce		json
// @Param		id		path		int	true	"actor id"
// @Success		200		{object}	statusResponse
// @Failure		400,403	{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/api/v1/actors/{id} [delete]
func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sid := mux.Vars(r)["id"]
	id, err := strconv.Atoi(sid)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//userRole := context.Get(r, USERROLE)
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		newErrorResponse(w, http.StatusForbidden, errors.New("not enough rights").Error())
		return
	}

	if err := h.service.Actor.DeleteActor(id); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	log.Println("DeleteActor with id - " + strconv.Itoa(id))
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(statusResponse{"ok"})
	_, _ = w.Write(result)
}
