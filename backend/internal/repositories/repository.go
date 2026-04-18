package repositories

import "github.com/jmoiron/sqlx"

type Repository struct {
	Users UserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUserRepo(db),
	}
}

type UserRepository interface {
	DeleteUserById(id int) error
}
