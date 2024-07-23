package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	GetCompany(id uint) (*model.Company, error)
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetCompany(id uint) (*model.Company, error) {
	company := &model.Company{}
	result := r.db.Find(company, "id = ?", id)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetCompany: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return company, nil
}
