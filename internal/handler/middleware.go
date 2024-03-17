package handler

import (
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

const (
	USERROLE = "userRole"
)

func (h *Handler) UserIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			newErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}
		headerPart := strings.Split(header, " ")
		if len(headerPart) != 2 {
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		userId, err := h.service.Authorization.ParseToken(headerPart[1])
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, "invalid auth header part")
			return
		}
		user, err := h.service.Authorization.GetUser(userId)
		if err != nil {
			newErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		context.Set(r, USERROLE, user.Role)
		next.ServeHTTP(w, r)
	})
}
