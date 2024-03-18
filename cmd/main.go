package main

import (
	"log"
	"vk/internal/app"
)

// @title			Swagger Film library
// @version		1.0
// @description	This is a sample server for Film library
// @host		localhost:8080
// @BasePath	/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
