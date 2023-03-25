package model

type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UnAuthRequest struct {
	Email string `json:"email" validate:"required"`
}
