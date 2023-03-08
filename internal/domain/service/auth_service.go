package service

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type AuthService interface {
	Authenticate(email, userPassword string) (authenticatedUser entity.User, err error)
	GenerateToken(user entity.User) (accessToken, refreshToken string, err error)
	RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error)
}
