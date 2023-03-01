package api

import (
	"github.com/devararishivian/antrekuy/internal/domain/service"
	"github.com/devararishivian/antrekuy/internal/presentation/model"
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
	req := new(model.StoreUserRequest)
	var res model.StoreUserResponse

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	//TODO: rewrite request validation and check valid role ID
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "all field must be filled")
	}

	storeRes, err := h.userService.Store(req.Name, req.Email, req.Password, req.RoleID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	res.Message = "success create user"

	//TODO: refactor create remapping method
	res.Data = model.UserResponse{
		ID:        storeRes.ID,
		Name:      storeRes.Name,
		Email:     storeRes.Email,
		CreatedAt: storeRes.CreatedAt,
		UpdatedAt: storeRes.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
