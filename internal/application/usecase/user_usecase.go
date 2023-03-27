package usecase

import (
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/repository"
	"github.com/devararishivian/go-clean-code-ddd/pkg/password"
	"github.com/devararishivian/go-clean-code-ddd/pkg/uuid"
)

type UserUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userRepository: userRepository,
	}
}

func (uc *UserUseCaseImpl) Store(name, email, reqPassword string, roleID int) (result *entity.User, err error) {
	newUser := new(entity.User)

	userID, err := uuid.NewUUID()
	if err != nil {
		return newUser, err
	}

	hashedPassword, err := password.HashPassword(reqPassword)
	if err != nil {
		return newUser, err
	}

	newUser.ID = userID
	newUser.Name = name
	newUser.Email = email
	newUser.Password = hashedPassword
	newUser.Role.ID = roleID

	result, err = uc.userRepository.Store(newUser)
	return
}

func (uc *UserUseCaseImpl) FindByEmail(email string) (entity.User, error) {
	user, err := uc.userRepository.FindByEmail(email)
	return user, err
}

func (uc *UserUseCaseImpl) FindByID(id string) (entity.User, error) {
	user, err := uc.userRepository.FindByID(id)
	return user, err
}
