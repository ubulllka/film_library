package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vk/internal/config"
)

const (
	USERS     = "users"
	ACTORS    = "actors"
	FILMS     = "films"
	FILMACTOR = "film_actor"
	PARSEDATE = "02-01-2006"
)

func runMigrations(db *sql.DB, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Error applying migrations: %v", err)
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				log.Fatalf("Error applying migrations: %v", err)
				return err
			}

			log.Printf("Applied migration: %s", file.Name())
		}
	}

	return nil
}

func InitializeDB() (*sql.DB, error) {
	info := config.GetConf().DB

	urlPostgres := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%d",
		info.User, info.Name, info.Password, info.Host, info.Port)

	db, err := sql.Open("postgres", urlPostgres)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		//log.Println(urlPostgres)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	if err := runMigrations(db, info.Path); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
		return nil, err
	}

	log.Print("Init database")
	return db, nil
}
