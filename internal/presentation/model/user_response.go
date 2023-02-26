package model

import "time"

type StoreUserResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"data"`
}
