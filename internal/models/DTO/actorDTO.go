package DTO

type ActorDTO struct {
	Name  string   `json:"name"`
	Sex   string   `json:"sex"`
	Date  string   `json:"date"`
	Films []string `json:"films"`
}
