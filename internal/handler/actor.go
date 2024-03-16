package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
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
	id := mux.Vars(r)["id"]

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
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
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

func UpdateActor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = id
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}
}

func DeleteActor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_ = id
	userRole := r.Context().Value(USERROLE)
	if userRole != "ADMIN" {
		http.Error(w, errors.New("not enough rights").Error(), http.StatusForbidden)
		return
	}
}
