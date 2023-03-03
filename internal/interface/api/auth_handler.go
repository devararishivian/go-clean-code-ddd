package api

import (
	appConfig "github.com/devararishivian/antrekuy/internal/config"
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/internal/interface/validator"
	"github.com/devararishivian/antrekuy/internal/presentation/model"
	"github.com/devararishivian/antrekuy/pkg/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validator
}

func NewAuthHandler(useCase service.AuthService) AuthHandler {
	return AuthHandler{
		authService: useCase,
		validator:   validator.NewValidator(),
	}
}

func (h *AuthHandler) Authenticate(c *fiber.Ctx) error {
	request := new(model.AuthRequest)
	var result model.AuthResponse

	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.validator.Validate(request); err.ValidationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	errAuthenticate := h.authService.Authenticate(request.Email, request.Password)
	if errAuthenticate != nil {
		return fiber.NewError(fiber.StatusUnauthorized, errAuthenticate.Error())
	}

	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	tokenExpiredTime := time.Now().Add(time.Hour * 24).Unix()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["refresh_token"] = refreshToken
	claims["exp"] = tokenExpiredTime

	t, err := token.SignedString([]byte(appConfig.JWTSecret))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	result.AccessToken = t
	result.RefreshToken = refreshToken
	result.ExpiresAt = tokenExpiredTime

	return c.Status(fiber.StatusOK).JSON(result)
}
