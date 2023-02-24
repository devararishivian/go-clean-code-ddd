package repository

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type UserRepository interface {
	Store(user *entity.User) error
}
