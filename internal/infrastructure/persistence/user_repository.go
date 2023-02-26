package persistence

import (
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
	return nil
}
