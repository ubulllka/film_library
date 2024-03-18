package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
	mock_handler "vk/internal/handler/mocks"
	"vk/internal/models/DTO"
)

func TestHandler_GetAllFilms(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockFilm, column, order string)

	arr := []DTO.FilmDTO{
		{
			Name:        "aname-film1",
			Description: "description-film",
			Date:        "10-10-0020",
			Rating:      9.8,
			Actors: []DTO.ActorInput{
				{
					Name: "name-actor1",
					Sex:  "male",
					Date: "01-01-0001",
				}, {
					Name: "name-actor2",
					Sex:  "female",
					Date: "02-02-0002",
				},
			},
		},
		{
			Name:        "bname-film2",
			Description: "description-film",
			Date:        "10-10-0050",
			Rating:      9.6,
			Actors:      make([]DTO.ActorInput, 0),
		},
	}
	byteArr, _ := json.Marshal(arr)
	tests := []struct {
		name                 string
		column               string
		order                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "Ok",
			column: "film_name",
			order:  "ASC",
			mockBehavior: func(r *mock_handler.MockFilm, column, order string) {
				r.EXPECT().GetAllFilms(column, order).Return(arr, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(byteArr),
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_handler.MockFilm, column, order string) {
				r.EXPECT().GetAllFilms(column, order).Return(nil, errors.New("something went wrong"))
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

			repo := mock_handler.NewMockFilm(c)
			test.mockBehavior(repo, test.column, test.order)

			services := &Service{Film: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/", handler.GetAllFilms).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf(`/?column=%s&order=%s`, test.column, test.order),
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_SearchFilms(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockFilm, fragment string)

	arr := []DTO.FilmDTO{
		{
			Name:        "aname-film1",
			Description: "description-film",
			Date:        "10-10-0020",
			Rating:      9.8,
			Actors: []DTO.ActorInput{
				{
					Name: "name-actor1",
					Sex:  "male",
					Date: "01-01-0001",
				}, {
					Name: "name-actor2",
					Sex:  "female",
					Date: "02-02-0002",
				},
			},
		},
		{
			Name:        "bname-film2",
			Description: "description-film",
			Date:        "10-10-0050",
			Rating:      9.6,
			Actors:      make([]DTO.ActorInput, 0),
		},
	}
	byteArr, _ := json.Marshal(arr)
	tests := []struct {
		name                 string
		fragment             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Ok",
			fragment: "a",
			mockBehavior: func(r *mock_handler.MockFilm, fragment string) {
				r.EXPECT().SearchFilms(fragment).Return(arr, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(byteArr),
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_handler.MockFilm, fragment string) {
				r.EXPECT().SearchFilms(fragment).Return(nil, errors.New("something went wrong"))
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

			repo := mock_handler.NewMockFilm(c)
			test.mockBehavior(repo, test.fragment)

			services := &Service{Film: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/", handler.SearchFilms).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf(`/?q=%s`, test.fragment),
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_GetFilm(t *testing.T) {
	type mockBehavior func(r *mock_handler.MockFilm, id int)

	film := DTO.FilmDTO{
		Name:        "aname-film1",
		Description: "description-film",
		Date:        "10-10-0020",
		Rating:      9.8,
		Actors: []DTO.ActorInput{
			{
				Name: "name-actor1",
				Sex:  "male",
				Date: "01-01-0001",
			}, {
				Name: "name-actor2",
				Sex:  "female",
				Date: "02-02-0002",
			},
		},
	}
	byteFilm, _ := json.Marshal(film)
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
			mockBehavior: func(r *mock_handler.MockFilm, id int) {
				r.EXPECT().GetFilm(id).Return(film, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(byteFilm),
		},
		{
			name:                 "Wrong Input",
			idURL:                "a",
			id:                   0,
			mockBehavior:         func(r *mock_handler.MockFilm, id int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
		},
		{
			name:  "Service Error",
			idURL: "1",
			id:    1,
			mockBehavior: func(r *mock_handler.MockFilm, id int) {
				r.EXPECT().GetFilm(id).Return(DTO.FilmDTO{}, errors.New("something went wrong"))
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

			repo := mock_handler.NewMockFilm(c)
			test.mockBehavior(repo, test.id)

			services := &Service{Film: repo}
			handler := NewHandler(services)

			// Init Endpoint
			r := mux.NewRouter()
			r.HandleFunc("/{id}", handler.GetFilm).Methods("GET")

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
