package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type CompanyUserRepository interface {
	GetCompanyUser(companyId, userId uint) (*model.CompanyUser, error)
	CreateCompanyUsers(companyUsers []*model.CompanyUser) ([]*model.CompanyUser, error)
}

type companyUserRepository struct {
	db *gorm.DB
}

func NewCompanyUserRepository(db *gorm.DB) CompanyUserRepository {
	return &companyUserRepository{db: db}
}

func (r *companyUserRepository) GetCompanyUser(companyId, userId uint) (*model.CompanyUser, error) {
	company := &model.CompanyUser{}
	result := r.db.Find(company, "company_id = ? AND user_id = ?", companyId, userId)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetCompanyUser: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return company, nil
}

func (r *companyUserRepository) CreateCompanyUsers(companyUsers []*model.CompanyUser) ([]*model.CompanyUser, error) {
	if err := r.db.Create(companyUsers).Error; err != nil {
		return nil, myErrors.ErrDb
	}
	return companyUsers, nil
}
