package models

import "time"

type Actor struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Sex  string    `json:"sex"`
	Date time.Time `json:"date"`
}
