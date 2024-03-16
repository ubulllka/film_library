package service

import (
	"vk/internal/models"
	"vk/internal/models/DTO"
	"vk/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GetUser(id int) (models.User, error)
}

type Actor interface {
	GetAllActors() ([]DTO.ActorDTO, error)
	GetActor(id string) (DTO.ActorDTO, error)
	CreateActor(actor models.Actor) (int, error)
}

type Film interface {
}

type Service struct {
	Authorization
	Actor
	Film
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Actor:         NewActorService(repo.Actor),
	}
}
