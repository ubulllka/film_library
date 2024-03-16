package app

import (
	"log"
	"net/http"
	"vk/internal/config"
	"vk/internal/db"
	"vk/internal/handler"
	"vk/internal/repository"
	"vk/internal/service"
)

func Run() error {

	if err := config.InitConfig(); err != nil {
		log.Panic(err)
		return err
	}
	conn, err := db.InitializeDB()
	if err != nil {
		log.Panic(err)
		return err
	}
	defer conn.Close()

	repo := repository.NewRepository(conn)
	serv := service.NewService(repo)
	hand := handler.NewHandler(serv)

	serverUrl := config.GetConf().Server.URL
	if err := http.ListenAndServe(serverUrl, hand.Routes()); err != nil {
		log.Fatalf("server did not start work: %s", err.Error())
		return err
	}
	return nil
}
