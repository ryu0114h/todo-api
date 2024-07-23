package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	result := r.db.Find(user, "username = ?", username)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetByUsername: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return user, nil
}

func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, myErrors.ErrDb
	}
	return user, nil
}
