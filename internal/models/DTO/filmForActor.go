package DTO

type FilmForActor struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Date        string  `json:"data"`
	Rating      float64 `json:"rating"`
}
