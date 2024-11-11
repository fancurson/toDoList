package repository

import (
	"fmt"
	"log"

	todo "github.com/fancurson/toDoList"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func newAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (p *AuthPostgres) CreateUser(user todo.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	var ID int
	err := p.db.QueryRow(query, user.Name, user.Username, user.Password).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User

	var id int

	log.Printf("u:%s, p:%s", username, password)
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", "users")
	err := r.db.QueryRowx(query, username, password).Scan(&id)
	log.Printf("r:%d", id)
	user.Id = id
	return user, err
}
