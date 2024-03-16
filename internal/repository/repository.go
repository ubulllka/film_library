package repository

import (
	"database/sql"
	"vk/internal/models"
	"vk/internal/models/DTO"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
}

type Actor interface {
	GetAllActors() ([]DTO.ActorDTO, error)
	getFilmNames(id string) ([]string, error)
	GetActor(id string) (DTO.ActorDTO, error)
	CreateActor(actor models.Actor) (int, error)
}

type Film interface {
}

type Repository struct {
	Authorization
	Actor
	Film
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Actor:         NewActorPostgres(db),
	}
}
