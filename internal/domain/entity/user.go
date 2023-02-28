package entity

import (
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []User

type UserRole struct {
	ID   int
	Name string
}
