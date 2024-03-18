package handler

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	mock_handler "vk/internal/handler/mocks"
)

func TestHandler_UserIdentify(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_handler.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_handler.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_handler.NewMockAuthorization(c)
			test.mockBehavior(repo, test.token)

			services := &Service{Authorization: repo}
			handler := Handler{services}

			// Init Endpoint
			r := mux.NewRouter()
			r.Use(handler.UserIdentify)
			r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				userRole := r.Context().Value(USERROLE)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"role": %s}`, userRole)))
			}).Methods("GET")

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
