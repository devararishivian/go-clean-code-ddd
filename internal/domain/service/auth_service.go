package service

import (
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Authenticate(email, userPassword string) (authenticatedUser entity.User, err error)
	GenerateToken(user entity.User) (accessToken, refreshToken string, err error)
	RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	ValidateToken(accessToken string) (claims jwt.MapClaims, errCode, errMessage string)
}
