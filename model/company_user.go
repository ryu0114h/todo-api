package model

import "time"

type CompanyUser struct {
	ID        uint
	CompanyID uint
	UserID    uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
