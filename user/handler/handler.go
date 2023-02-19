package handler

import "github.com/devararishivian/antrekuy/user/domain"

type Handler struct {
	UseCase domain.UserUseCase
}

func NewHandler(useCase domain.UserUseCase) Handler {
	return Handler{
		useCase,
	}
}
