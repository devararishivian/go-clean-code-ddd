package usecase

import (
	"github.com/devararishivian/antrekuy/user/domain"
)

type UseCase struct {
	repo domain.UserRepository
}

func NewUseCase(repo domain.UserRepository) *UseCase {
	return &UseCase{repo}
}

func (s *UseCase) Store(req *domain.NewUserRequest) error {
	return nil
}
