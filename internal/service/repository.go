package service

import (
	"database/sql"
	"vk/internal/models"
	"vk/internal/models/DTO"
	"vk/internal/repository"
)

type Authorization interface {
	Create(user models.User) (int, error)
	GetOne(username, password string) (models.User, error)
	GetOneById(id int) (models.User, error)
}

type Actor interface {
	GetAll() ([]DTO.ActorDTO, error)
	GetOne(id int) (DTO.ActorDTO, error)
	Create(actor models.Actor) (int, error)
	Update(id int, input DTO.ActorUpdate) error
	Delete(id int) error
}

type Film interface {
	GetAll(column, order string) ([]DTO.FilmDTO, error)
	GetOne(id int) (DTO.FilmDTO, error)
	Search(fragment string) ([]DTO.FilmDTO, error)
	Create(film models.Film, arr []int) (int, error)
	Update(id int, input DTO.FilmUpdate) error
	Delete(id int) error
}

type Repository struct {
	Authorization
	Actor
	Film
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: repository.NewAuthPostgres(db),
		Actor:         repository.NewActorPostgres(db),
		Film:          repository.NewFilmPostgres(db),
	}
}
