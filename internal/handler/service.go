package handler

import (
	"vk/internal/models"
	"vk/internal/models/DTO"
	"vk/internal/service"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GetUser(id int) (models.User, error)
}

type Actor interface {
	GetAllActors() ([]DTO.ActorDTO, error)
	GetActor(id int) (DTO.ActorDTO, error)
	CreateActor(actor models.Actor) (int, error)
	UpdateActor(id int, input DTO.ActorUpdate) error
	DeleteActor(id int) error
}

type Film interface {
	GetAllFilms(column, order string) ([]DTO.FilmDTO, error)
	GetFilm(id int) (DTO.FilmDTO, error)
	SearchFilms(fragment string) ([]DTO.FilmDTO, error)
	CreateFilm(film models.Film, arr []int) (int, error)
	UpdateFilm(id int, input DTO.FilmUpdate) error
	DeleteFilm(id int) error
}

type Service struct {
	Authorization
	Actor
	Film
}

func NewService(repo *service.Repository) *Service {
	return &Service{
		Authorization: service.NewAuthService(repo.Authorization),
		Actor:         service.NewActorService(repo.Actor),
		Film:          service.NewFilmService(repo.Film),
	}
}
