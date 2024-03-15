package models

import "time"

type Film struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"data"`
	Rating      float64   `json:"rating"`
}

type FilmActor struct {
	Id      int
	FilmId  int
	ActorId int
}
