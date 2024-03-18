package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	mock_handler "vk/internal/handler/mocks"
	"vk/internal/models/DTO"
)

func TestHandler_GetAllActors(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockActor)

	arr := []DTO.ActorDTO{
		{
			Name:  "name-actor1",
			Sex:   "male",
			Date:  "01-01-0001",
			Films: make([]DTO.FilmForActor, 0),
		}, {
			Name: "name-actor2",
			Sex:  "female",
			Date: "02-02-0002",
			Films: []DTO.FilmForActor{
				{
					Name:        "name-film",
					Description: "description-film",
					Date:        "10-10-0050",
					Rating:      9.8,
				},
			},
		},
	}
	byteArr, _ := json.Marshal(arr)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mock_handler.MockActor) {
				r.EXPECT().GetAllActors().Return(arr, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(byteArr),
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_handler.MockActor) {
				r.EXPECT().GetAllActors().Return(nil, errors.New("something went wrong"))
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

			repo := mock_handler.NewMockActor(c)
			test.mockBehavior(repo)

			services := &Service{Actor: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/", handler.GetAllActors).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/",
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_GetActor(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockActor, id int)

	actor := DTO.ActorDTO{
		Name: "name-actor2",
		Sex:  "female",
		Date: "02-02-0002",
		Films: []DTO.FilmForActor{
			{
				Name:        "name-film",
				Description: "description-film",
				Date:        "10-10-0050",
				Rating:      9.8,
			},
		},
	}
	byteActor, _ := json.Marshal(actor)
	tests := []struct {
		name                 string
		idURL                string
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "Ok",
			idURL: "1",
			id:    1,
			mockBehavior: func(r *mock_handler.MockActor, id int) {
				r.EXPECT().GetActor(id).Return(actor, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(byteActor),
		},
		{
			name:                 "Wrong Input",
			idURL:                "a",
			id:                   0,
			mockBehavior:         func(r *mock_handler.MockActor, id int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
		},
		{
			name:  "Service Error",
			idURL: "1",
			id:    1,
			mockBehavior: func(r *mock_handler.MockActor, id int) {
				r.EXPECT().GetActor(id).Return(DTO.ActorDTO{}, errors.New("something went wrong"))
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

			repo := mock_handler.NewMockActor(c)
			test.mockBehavior(repo, test.id)

			services := &Service{Actor: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/{id}", handler.GetActor).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/"+test.idURL,
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
