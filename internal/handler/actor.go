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

func (h *Handler) GetALLActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.Actor.GetAllActors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actorsByte, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorsByte)
}

func (h *Handler) GetActor(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["id"]

	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actor, err := h.service.Actor.GetActor(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actorByte, err := json.Marshal(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(actorByte)

}

func (h *Handler) SaveActor(w http.ResponseWriter, r *http.Request) {
	userRole := context.Get(r, USERROLE)
	fmt.Println(userRole)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}

	var actorSave DTO.ActorInput
	if err := json.NewDecoder(r.Body).Decode(&actorSave); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actor, err := actorSave.GetActor()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(actor); err != nil {
		errs := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errs), http.StatusBadRequest)
		return
	}

	id, err := h.service.Actor.CreateActor(actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Actor %d create\"}", id))))

}

func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
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

	var input DTO.ActorUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Actor.UpdateActor(id, input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Actor %d update\"}", id))))
}

func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.Actor.DeleteActor(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf("{ \"message\": \"Actor %d delete\"}", id))))
}
