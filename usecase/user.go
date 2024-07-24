package usecase

import (
	"todo-api/model"
	"todo-api/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	CreateUser(username, email, password, role string, companyIds []uint) (*model.User, error)
}

type userUseCase struct {
	db                    *gorm.DB
	userRepository        repository.UserRepository
	companyUserRepository repository.CompanyUserRepository
}

func NewUserUseCase(
	db *gorm.DB,
	userRepository repository.UserRepository,
	companyUserRepository repository.CompanyUserRepository,
) UserUseCase {
	return &userUseCase{
		db:                    db,
		userRepository:        userRepository,
		companyUserRepository: companyUserRepository,
	}
}

func (u *userUseCase) CreateUser(username, email, password, role string, companyIds []uint) (*model.User, error) {

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// トランザクションの開始
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}
	createdUser, err := u.userRepository.CreateUser(user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	companyUsers := []*model.CompanyUser{}
	for _, companyId := range companyIds {
		companyUsers = append(companyUsers, &model.CompanyUser{
			UserID:    createdUser.ID,
			CompanyID: companyId,
		})
	}
	_, err = u.companyUserRepository.CreateCompanyUsers(companyUsers)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// トランザクションのコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return createdUser, nil
}
