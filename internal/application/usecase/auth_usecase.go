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
	"strings"
	"time"
)

// TODO: Refactor error const

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

func (au *AuthUseCaseImpl) Authenticate(email, userPassword string) (authenticatedUser entity.User, err error) {
	user, err := au.userUseCase.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		return user, errors.New("no user with given email/password")
	}

	errComparePassword := password.ComparePassword(user.Password, userPassword)
	if errComparePassword != nil {
		return user, errors.New("no user with given email/password")
	}

	return user, nil
}

func (au *AuthUseCaseImpl) GenerateToken(user entity.User) (accessToken, refreshToken string, err error) {
	accessToken, err = au.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = au.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = au.storeTokenToCache(user.ID, accessToken, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (au *AuthUseCaseImpl) generateAccessToken(user entity.User) (accessToken string, err error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role.ID,
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func (au *AuthUseCaseImpl) RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	tokenClaims, errCode, errMessage := au.ValidateToken(accessToken)
	if errCode != "" {
		if errCode != "ErrExpiredToken" {
			return "", "", fmt.Errorf("failed to validate access token: %v", errMessage)
		}
	}

	userIDFromClaims := tokenClaims["id"].(string)
	_, cachedRefreshToken, err := au.getTokenFromCache(userIDFromClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to get cached token: %v", err)
	}

	if cachedRefreshToken != refreshToken {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := au.userUseCase.FindByID(userIDFromClaims)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %v", err)
	}

	newAccessToken, newRefreshToken, err = au.GenerateToken(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new token: %v", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (au *AuthUseCaseImpl) generateRefreshToken() (refreshToken string, err error) {
	refreshToken, err = uuid.NewUUID()
	if err != nil {
		return "", errors.New("failed to generate refresh token")
	}

	return refreshToken, nil
}

func (au *AuthUseCaseImpl) ValidateToken(accessToken string) (claims jwt.MapClaims, errCode, errMessage string) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(appConfig.JWTSecret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, "ErrSignatureInvalid", jwt.ErrSignatureInvalid.Error()
		} else if err.Error() != fmt.Sprintf("%s: %s", jwt.ErrTokenInvalidClaims.Error(), jwt.ErrTokenExpired.Error()) {
			return nil, "ErrInvalidToken", err.Error()
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		expirationTime, ok := claims["exp"].(float64)
		if !ok {
			return nil, "ErrInvalidToken", jwt.ErrTokenInvalidClaims.Error()
		}

		if time.Now().Unix() > int64(expirationTime) {
			return claims, "ErrExpiredToken", jwt.ErrTokenExpired.Error()
		}

		return nil, "ErrInvalidToken", jwt.ErrTokenInvalidClaims.Error()
	}

	return claims, "", ""
}

func (au *AuthUseCaseImpl) storeTokenToCache(userID, accessToken, refreshToken string) error {
	formattedUserID := strings.ReplaceAll(userID, "-", "")

	cacheData := entity.Cache{
		Key: fmt.Sprintf("auth:%s", formattedUserID),
		Value: map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
		TTL: time.Hour * 24 * 7,
	}

	err := au.cacheRepository.HSet(cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (au *AuthUseCaseImpl) getTokenFromCache(userID string) (accessToken, refreshToken string, err error) {
	formattedUserID := strings.ReplaceAll(userID, "-", "")

	val, err := au.cacheRepository.HGetAll(fmt.Sprintf("auth:%s", formattedUserID))
	if err != nil {
		return "", "", err
	}

	valMap := make(map[string]string, len(val.Value.(map[string]any)))
	for k, v := range val.Value.(map[string]any) {
		valMap[k] = v.(string)
	}

	return valMap["access_token"], valMap["refresh_token"], nil
}
