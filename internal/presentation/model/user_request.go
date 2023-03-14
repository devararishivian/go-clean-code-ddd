package model

type StoreUserRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,min=6,max=12"`
	RoleID   int    `json:"role_id" validate:"required,oneof=1 2 3 4"`
}
