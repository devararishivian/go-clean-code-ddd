package repository

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type UserRepository interface {
	Store(user *entity.User) (*entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindByID(id string) (entity.User, error)
}
