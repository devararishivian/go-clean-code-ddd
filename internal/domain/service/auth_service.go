package service

import (
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Authenticate(email, userPassword string) (authenticatedUser entity.User, err error)
	UnAuthenticate(email string) error
	GenerateToken(user entity.User) (accessToken string, err error)
	ValidateToken(accessToken string) (claims jwt.MapClaims, errCode, errMessage string)
}
