package app

import (
	"log"
	"net/http"
	"vk/internal/config"
	"vk/internal/db"
	"vk/internal/handler"
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

	repo := service.NewRepository(conn)
	serv := handler.NewService(repo)
	hand := handler.NewHandler(serv)

	log.Println("Init repositories, services, handlers")

	serverUrl := config.GetConf().Server.URL
	if err := http.ListenAndServe(serverUrl, hand.Routes()); err != nil {
		log.Fatalf("server did not start work: %s", err.Error())
		return err
	}
	log.Println("Server listening url " + serverUrl)
	return nil
}
