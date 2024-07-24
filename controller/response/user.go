package response

import (
	"time"
	"todo-api/model"
)

type CreateUserResponseBody struct {
	User *CreateUserResponseBodyUser `json:"user"`
}

type CreateUserResponseBodyUser struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         string     `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func NewCreateUserResponseBody(user *model.User) *CreateUserResponseBody {
	return &CreateUserResponseBody{
		User: &CreateUserResponseBodyUser{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		},
	}
}
