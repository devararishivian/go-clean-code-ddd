package service

import "github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"

type UserService interface {
	Store(name, email, password string, roleID int) (*entity.User, error)
	FindByEmail(email string) (entity.User, error)
	FindByID(id string) (entity.User, error)
}
