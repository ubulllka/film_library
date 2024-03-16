package models

import "time"

type Actor struct {
	Id   int       `json:"-"`
	Name string    `json:"name" validate:"required"`
	Sex  string    `json:"sex" validate:"required"`
	Date time.Time `json:"date" validate:"required"`
}
