package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vk/internal/config"
)

func runMigrations(db *pgx.Conn, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Error applying migrations: %v", err)
				return err
			}

			_, err = db.Exec(context.Background(), string(content))
			if err != nil {
				log.Fatalf("Error applying migrations: %v", err)
				return err
			}

			log.Printf("Applied migration: %s", file.Name())
		}
	}

	return nil
}

func InitializeDB() error {
	info := config.CONFIG.DB
	urlPostgres := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		info.User, info.Password, info.Host, info.Port, info.Name)

	conn, err := pgx.Connect(context.Background(), urlPostgres)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		log.Println(urlPostgres)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Panic(err)
		return err
	}

	if err := runMigrations(conn, info.Path); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	log.Print("Init database")
	return nil
}
