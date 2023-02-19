package repository

import (
	"database/sql"
	"github.com/devararishivian/antrekuy/user/domain"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (u *Repository) Store(req *domain.NewUserRequest) error {
	return nil
}
