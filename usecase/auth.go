package usecase

import (
	"errors"
	"todo-api/model"
	"todo-api/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(username, password string) (*model.User, error)
}

type authUseCase struct {
	userRepository repository.UserRepository
}

func NewAuthUseCase(userRepository repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepository: userRepository}
}

func (a *authUseCase) Login(username, password string) (*model.User, error) {
	user, err := a.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
