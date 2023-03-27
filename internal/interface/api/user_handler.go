package api

import (
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/service"
	"github.com/devararishivian/go-clean-code-ddd/internal/interface/validator"
	"github.com/devararishivian/go-clean-code-ddd/internal/presentation/model"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
}

func NewUserHandler(useCase service.UserService) UserHandler {
	return UserHandler{
		userService: useCase,
		validator:   validator.NewValidator(),
	}
}

func (h *UserHandler) Store(c *fiber.Ctx) error {
	req := new(model.StoreUserRequest)
	var res model.StoreUserResponse

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := h.validator.Validate(req); err.ValidationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
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
