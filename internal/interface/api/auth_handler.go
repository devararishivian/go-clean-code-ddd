package api

import (
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/internal/interface/validator"
	"github.com/devararishivian/antrekuy/internal/presentation/model"
	"github.com/gofiber/fiber/v2"
	"strings"
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

	authenticatedUser, err := h.authService.Authenticate(request.Email, request.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	accessToken, refreshToken, err := h.authService.GenerateToken(authenticatedUser)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	result.AccessToken = accessToken
	result.RefreshToken = refreshToken

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	request := new(model.RefreshTokenRequest)
	var result model.AuthResponse

	if err := c.BodyParser(request); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.validator.Validate(request); err.ValidationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	authHeader := c.Get("Authorization")

	// Split the Authorization header into two authHeaderParts: Bearer and the token
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid Authorization header format",
		})
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(authHeaderParts[1], request.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	result.AccessToken = newAccessToken
	result.RefreshToken = newRefreshToken
	return c.Status(fiber.StatusOK).JSON(result)
}
