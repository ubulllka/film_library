package handler

import (
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

const (
	USERID   = "userID"
	USERROLE = "userRole"
)

func (h *Handler) UserIdentify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, fmt.Sprintf("empty auth header"), http.StatusUnauthorized)
			return
		}
		headerPart := strings.Split(header, " ")
		if len(headerPart) != 2 {
			http.Error(w, fmt.Sprintf("invalid auth header"), http.StatusUnauthorized)
			return
		}

		userId, err := h.service.Authorization.ParseToken(headerPart[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid auth header part"), http.StatusUnauthorized)
			return
		}
		user, err := h.service.Authorization.GetUser(userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		fmt.Println(user.Id, user.Username, user.Role, user.Password)

		context.Set(r, USERID, userId)
		context.Set(r, USERROLE, user.Role)
		next.ServeHTTP(w, r)
	})
}
