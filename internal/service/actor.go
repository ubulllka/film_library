package service

import (
	"vk/internal/models"
	"vk/internal/models/DTO"
	"vk/internal/repository"
)

type ActorService struct {
	repo repository.Actor
}

func NewActorService(repo repository.Actor) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) GetAllActors() ([]DTO.ActorDTO, error) {
	return s.repo.GetAllActors()
}

func (s *ActorService) GetActor(id string) (DTO.ActorDTO, error) {
	return s.repo.GetActor(id)
}
func (s *ActorService) CreateActor(actor models.Actor) (int, error) {
	return s.repo.CreateActor(actor)
}
