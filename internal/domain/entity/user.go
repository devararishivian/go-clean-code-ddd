package entity

import (
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []User

type UserRole struct {
	ID   int
	Name string
}

type UserHasRole struct {
	UserID string
	RoleID int
}

type NewUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"roleID"`
}