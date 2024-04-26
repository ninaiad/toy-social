package database

import (
	"soa-main/internal/user"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user user.User) error
	GetUser(username, password string) (user.User, error)
	GetUserData(userId int) (user.UserPublic, error)
	UpdateUser(userId int, update user.UserPublic, timeUpdated string) error
}

type Database struct {
	Authorization
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		Authorization: NewAuthPostgres(db),
	}
}
