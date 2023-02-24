package service

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type UserService interface {
	Store(name, email, password string) (*entity.User, error)
}
