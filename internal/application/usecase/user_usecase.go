package usecase

import (
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/pkg/uuid"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepository: userRepository,
	}
}

func (uc *UserUseCaseImpl) Store(name, email, password string) (*entity.User, error) {
	newUser := new(entity.User)

	userID, err := uuid.NewUUID()
	if err != nil {
		return newUser, err
	}

	newUser.ID = userID
	newUser.Name = name
	newUser.Email = email
	newUser.Password = password

	if err := uc.userRepository.Store(newUser); err != nil {
		return newUser, err
	}

	return newUser, nil
}
