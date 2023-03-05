package usecase

import (
	"errors"
	"fmt"
	appConfig "github.com/devararishivian/antrekuy/internal/config"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/pkg/password"
	"github.com/devararishivian/antrekuy/pkg/uuid"
	"github.com/golang-jwt/jwt/v5"
	"time"
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

func (s *AuthUseCaseImpl) GenerateToken(email string) (accessToken, refreshToken string, err error) {
	// Generate the access token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	accessToken, err = token.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	// Generate the refresh token
	refreshToken, err = uuid.NewUUID()
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	cacheData := entity.Cache{
		Key:   fmt.Sprintf("refresh_token:%s", email),
		Value: refreshToken,
		TTL:   time.Hour * 24 * 7,
	}

	err = s.cacheRepository.Set(cacheData)
	if err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}
