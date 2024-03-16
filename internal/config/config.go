package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

var CONFIG Config

type Config struct {
	Env    string `yaml:"env"`
	Server struct {
		URL string `yaml:"address" env-default:"localhost:8080"`
	} `yaml:"server"`
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Path     string `yaml:"path"`
	} `yaml:"db"`
}

func InitConfig() error {
	if err := cleanenv.ReadConfig("./config/local.yml", &CONFIG); err != nil {
		log.Fatal(err)
		return err
	}
	log.Print("Init config")
	return nil
}

func GetConf() Config {
	return CONFIG
}
