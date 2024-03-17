package service

import (
	"vk/internal/models"
	"vk/internal/models/DTO"
)

type FilmService struct {
	repo Film
}

func NewFilmService(repo Film) *FilmService {
	return &FilmService{repo: repo}
}

func (s *FilmService) GetAllFilms(column, order string) ([]DTO.FilmDTO, error) {
	return s.repo.GetAll(column, order)
}

func (s *FilmService) GetFilm(id int) (DTO.FilmDTO, error) {
	return s.repo.GetOne(id)
}

func (s *FilmService) SearchFilms(fragment string) ([]DTO.FilmDTO, error) {
	return s.repo.Search(fragment)
}

func (s *FilmService) CreateFilm(film models.Film, arr []int) (int, error) {
	return s.repo.Create(film, arr)
}

func (s *FilmService) UpdateFilm(id int, input DTO.FilmUpdate) error {
	return s.repo.Update(id, input)
}

func (s *FilmService) DeleteFilm(id int) error {
	return s.repo.Delete(id)
}
