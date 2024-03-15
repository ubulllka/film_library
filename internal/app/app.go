package app

import (
	"log"
	"net/http"
	"vk/internal/config"
	"vk/internal/db"
	"vk/internal/handler"
)

func Run() error {

	if err := config.InitConfig(); err != nil {
		log.Panic(err)
		return err
	}

	if err := db.InitializeDB(); err != nil {
		log.Panic(err)
		return err
	}

	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		log.Fatalf("server did not start work: %s", err.Error())
		return err
	}
	return nil
}
