package DTO

type FilmDTO struct {
	Id          int          `json:"-"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Date        string       `json:"data"`
	Rating      float64      `json:"rating"`
	Actors      []ActorInput `json:"actors"`
}
