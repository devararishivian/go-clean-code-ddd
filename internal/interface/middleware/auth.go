package middleware

import (
	"fmt"
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

const (
	authHeader      = "Authorization"
	authBearerToken = "Bearer "
)

type AuthMiddleware struct {
	authUseCase service.AuthService
}

func NewAuthMiddleware(authUseCase service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

func (a *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authHeader)

		if authorizationHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    "ErrMissingAuthHeader",
				"message": "Authorization header is missing",
			})
		}

		if !strings.HasPrefix(authorizationHeader, authBearerToken) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    "ErrInvalidAuthHeaderFormat",
				"message": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authorizationHeader, authBearerToken)
		tokenClaims, errCode, errMessage := a.authUseCase.ValidateToken(tokenString)
		if errCode != "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    errCode,
				"message": errMessage,
			})
		}

		userID, ok := tokenClaims["id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
		}

		roleID, ok := tokenClaims["role"]
		if !ok {
			fmt.Println(tokenClaims["role"])
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
		}

		// Set the user data in the request context
		c.Locals("userID", userID)
		c.Locals("roleID", roleID)

		// Call the next middleware
		return c.Next()
	}
}
