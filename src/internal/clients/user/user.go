package userClient

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"user_name"`
	Email    string    `json:"email"`
}
