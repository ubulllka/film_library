package DTO

import (
	"time"
	"vk/internal/db"
	"vk/internal/models"
)

type ActorInput struct {
	Name string `json:"name"`
	Sex  string `json:"sex"`
	Date string `json:"date"`
}

func (a *ActorInput) GetActor() (models.Actor, error) {
	t, err := time.Parse(db.PARSEDATE, a.Date)
	if err != nil {
		return models.Actor{}, err
	}
	return models.Actor{
		Name: a.Name,
		Sex:  a.Sex,
		Date: t,
	}, nil
}
