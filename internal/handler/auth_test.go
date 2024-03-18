package handler

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	mock_handler "vk/internal/handler/mocks"
	"vk/internal/models"
)

func TestHandler_SingUp(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockAuthorization, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty", "role": "USER"}`,
			inputUser: models.User{
				Username: "username",
				Password: "qwerty",
				Role:     "USER",
			},
			mockBehavior: func(r *mock_handler.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id": 1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            models.User{},
			mockBehavior:         func(r *mock_handler.MockAuthorization, user models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'User.Role' Error:Field validation for 'Role' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty", "role": "USER"}`,
			inputUser: models.User{
				Username: "username",
				Password: "qwerty",
				Role:     "USER",
			},
			mockBehavior: func(r *mock_handler.MockAuthorization, user models.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_handler.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &Service{Authorization: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/sing-up", handler.SingUp).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sing-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_SingIn(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockAuthorization, user UserInType)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            UserInType
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: UserInType{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_handler.MockAuthorization, user UserInType) {
				r.EXPECT().GenerateToken(user.Username, user.Password).Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            UserInType{},
			mockBehavior:         func(r *mock_handler.MockAuthorization, user UserInType) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'UserInType.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputUser: UserInType{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_handler.MockAuthorization, user UserInType) {
				r.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_handler.NewMockAuthorization(c)
			test.mockBehavior(repo, test.inputUser)

			services := &Service{Authorization: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/sing-in", handler.SingIn).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sing-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
