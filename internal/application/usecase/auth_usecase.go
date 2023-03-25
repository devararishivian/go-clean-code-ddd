package usecase

import (
	"errors"
	"fmt"
	appConfig "github.com/devararishivian/antrekuy/internal/config"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/pkg/password"
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

func (au *AuthUseCaseImpl) GenerateToken(user entity.User) (accessToken string, err error) {
	accessToken, err = au.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	err = au.storeTokenToCache(user.ID, accessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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

	existingAccessToken, err := au.getTokenFromCache(claims["id"].(string))
	if err != nil {
		return nil, "ErrInvalidToken", err.Error()
	}

	if existingAccessToken != accessToken {
		return nil, "ErrInvalidToken", jwt.ErrTokenNotValidYet.Error()
	}

	return claims, "", ""
}

func (au *AuthUseCaseImpl) storeTokenToCache(userID, accessToken string) error {
	formattedUserID := strings.ReplaceAll(userID, "-", "")

	cacheData := entity.Cache{
		Key: fmt.Sprintf("auth:%s", formattedUserID),
		Value: map[string]interface{}{
			"access_token": accessToken,
		},
		TTL: time.Hour * 24 * 7,
	}

	err := au.cacheRepository.HSet(cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (au *AuthUseCaseImpl) getTokenFromCache(userID string) (accessToken string, err error) {
	formattedUserID := strings.ReplaceAll(userID, "-", "")
	key := "auth:" + formattedUserID

	val, err := au.cacheRepository.HGetAll(key)
	if err != nil {
		return "", err
	}

	valMap := map[string]string{}
	if val.Value != nil {
		if valMapRaw, ok := val.Value.(map[string]any); ok {
			for k, v := range valMapRaw {
				valMap[k] = v.(string)
			}
		} else {
			return "", fmt.Errorf("unexpected value type: %T", val.Value)
		}
	}

	return valMap["access_token"], nil
}
