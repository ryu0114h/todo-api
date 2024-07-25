package model

import "time"

type Task struct {
	ID           uint
	CompanyID    uint
	CreateUserId uint
	Title        string
	Description  string
	DueDate      *time.Time
	AssigneeID   *uint
	Visibility   string
	Status       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	Assignee     *User
}
