package models

import "time"

type Film struct {
	Id          int       `json:"-"`
	Name        string    `json:"name" validate:"required,min=1,max=150"`
	Description string    `json:"description" validate:"required,lte=1000"`
	Date        time.Time `json:"data" validate:"required"`
	Rating      float64   `json:"rating" validate:"required,min=0,max=10"`
}

type FilmActor struct {
	Id      int
	FilmId  int
	ActorId int
}
