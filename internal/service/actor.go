package service

import (
	"vk/internal/models"
	"vk/internal/models/DTO"
)

type ActorService struct {
	repo Actor
}

func NewActorService(repo Actor) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) GetAllActors() ([]DTO.ActorDTO, error) {
	return s.repo.GetAll()
}

func (s *ActorService) GetActor(id int) (DTO.ActorDTO, error) {
	return s.repo.GetOne(id)
}

func (s *ActorService) CreateActor(actor models.Actor) (int, error) {
	return s.repo.Create(actor)
}

func (s *ActorService) UpdateActor(id int, input DTO.ActorUpdate) error {
	return s.repo.Update(id, input)
}

func (s *ActorService) DeleteActor(id int) error {
	return s.repo.Delete(id)
}
