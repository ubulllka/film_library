package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"vk/internal/models"
)

// SingUp
// @Summary		Sing Up
// @Tags		auth
// @Description	create account
// @Accept		json
// @Produce		json
// @Param		user	body		models.User	true	"account info"
// @Success		200		{integer}	integer		1
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/auth/sign-up [post]
func (h *Handler) SingUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		errs := err.(validator.ValidationErrors)
		newErrorResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	id, err := h.service.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("SingUp is ok")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

type UserInType struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SingIn
// @Summary		Sing In
// @Tags		auth
// @Description	login
// @Accept		json
// @Produce		json
// @Param		user	body		UserInType	true	"credentials"
// @Success		200		{string}	string		"token"
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router		/auth/sign-in [post]
func (h *Handler) SingIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user UserInType

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	validate := validator.New()

	if err := validate.Struct(user); err != nil {
		errs := err.(validator.ValidationErrors)
		newErrorResponse(w, http.StatusBadRequest, errs.Error())
		return
	}

	token, err := h.service.Authorization.GenerateToken(user.Username, user.Password)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("SingIn is ok")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"token\": \"%s\"}", token)))
}
