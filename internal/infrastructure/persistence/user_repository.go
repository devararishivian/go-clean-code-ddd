package persistence

import (
	"context"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
)

type UserRepositoryImpl struct {
	db *infrastructure.Database
}

func NewUserRepository(db *infrastructure.Database) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) Store(user *entity.User) error {
	const sql string = `INSERT INTO "user"(id, name, email, password) VALUES ($1, $2, $3, $4)`

	_, err := u.db.Conn.Exec(
		context.Background(),
		sql, user.ID,
		user.Name,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}
