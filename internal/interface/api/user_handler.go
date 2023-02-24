package api

import (
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(useCase service.UserService) UserHandler {
	return UserHandler{
		useCase,
	}
}

func (h *UserHandler) Store(c *fiber.Ctx) error {
	res, err := h.userService.Store("nama", "emil", "paswot")
	if err != nil {
		return c.JSON("error")
	}

	return c.JSON(res)
}
