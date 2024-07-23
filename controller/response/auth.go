package response

import (
	"time"
	"todo-api/model"
)

type LoginResponseBody struct {
	User *LoginResponseBodyUser `json:"user"`
}

type LoginResponseBodyUser struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         string     `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func NewLoginResponseBody(user *model.User) *LoginResponseBody {
	return &LoginResponseBody{
		User: &LoginResponseBodyUser{
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
