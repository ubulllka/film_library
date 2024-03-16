package repository

import (
	"database/sql"
	"fmt"
	"vk/internal/db"
	"vk/internal/models"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, user_role) values ($1, $2, $3) RETURNING id", db.USERS)

	row := r.db.QueryRow(query, user.Username, user.Password, user.Role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var id int
	var role string
	query := fmt.Sprintf("SELECT id, user_role FROM %s WHERE username=$1 AND password_hash=$2", db.USERS)
	if err := r.db.QueryRow(query, username, password).Scan(&id, &role); err != nil {
		return models.User{}, err
	}
	return models.User{Id: id, Username: username, Password: password, Role: role}, nil
}

func (r *AuthPostgres) GetUserById(id int) (models.User, error) {
	var username, password, role string
	query := fmt.Sprintf("SELECT username, password_hash, user_role FROM %s WHERE id=$1", db.USERS)
	if err := r.db.QueryRow(query, id).Scan(&username, &password, &role); err != nil {
		return models.User{}, err
	}
	return models.User{Id: id, Username: username, Password: password, Role: role}, nil
}
