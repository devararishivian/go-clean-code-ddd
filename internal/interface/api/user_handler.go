package api

import (
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/internal/presentation/model"
	"github.com/gofiber/fiber/v2"
	"time"
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
	var res model.StoreUserResponse

	storeRes, err := h.userService.Store("nama", "emil", "paswot")
	if err != nil {
		res.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res.Message = "success create user"
	res.Data = struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}(struct {
		ID        string
		Name      string
		Email     string
		CreatedAt time.Time
		UpdatedAt time.Time
	}{
		ID:        storeRes.ID,
		Name:      storeRes.Name,
		Email:     storeRes.Email,
		CreatedAt: storeRes.CreatedAt,
		UpdatedAt: storeRes.UpdatedAt,
	})

	return c.Status(fiber.StatusOK).JSON(res)
}
