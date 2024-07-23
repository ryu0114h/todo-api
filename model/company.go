package model

import "time"

type Company struct {
	ID        uint
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
