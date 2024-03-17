package DTO

import (
	"errors"
	"time"
	"vk/internal/db"
	"vk/internal/models"
)

type FilmInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"data"`
	Rating      float64 `json:"rating"`
	Actors      []int   `json:"actors"`
}

func (f *FilmInput) GetFilmAndActors() (models.Film, []int, error) {
	t, err := time.Parse(db.PARSEDATE, f.Date)
	if err != nil {
		return models.Film{}, nil, err
	}
	if f.Actors == nil {
		return models.Film{}, nil, errors.New("actors list is null")
	}
	return models.Film{
		Name:        f.Name,
		Description: f.Description,
		Date:        t,
		Rating:      f.Rating,
	}, f.Actors, nil
}
