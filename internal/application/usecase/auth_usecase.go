package usecase

import (
	"errors"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/pkg/password"
)

type AuthUseCaseImpl struct {
	userUseCase     service.UserService
	cacheRepository repository.CacheRepository
}

func NewAuthUseCase(userUseCase service.UserService, cacheRepository repository.CacheRepository) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		userUseCase,
		cacheRepository,
	}
}

func (au *AuthUseCaseImpl) Authenticate(email, userPassword string) (err error) {
	user, err := au.userUseCase.FindByEmail(email)
	if err != nil {
		return err
	}

	if user.ID == "" {
		return errors.New("no user with given email/password")
	}

	errComparePassword := password.ComparePassword(user.Password, userPassword)
	if errComparePassword != nil {
		return errors.New("no user with given email/password")
	}

	return nil
}

func (au *AuthUseCaseImpl) Revoke(email string) error {
	return nil
}
