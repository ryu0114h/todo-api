package request

import (
	"time"
	"todo-api/model"
)

type CreateTaskRequestBody struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uint      `json:"assignee_id"`
	Visibility  string     `json:"visibility" validate:"required,oneof=company private"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress done"`
}

func NewTaskFromCreateTaskRequestBody(companyId uint, createUserId uint, requestBody *CreateTaskRequestBody) *model.Task {
	return &model.Task{
		CompanyID:    companyId,
		CreateUserId: createUserId,
		Title:        requestBody.Title,
		Description:  requestBody.Description,
		DueDate:      requestBody.DueDate,
		AssigneeID:   requestBody.AssigneeID,
		Visibility:   requestBody.Visibility,
		Status:       requestBody.Status,
	}
}

type CreateTaskByAdminRequestBody struct {
	CompanyId   uint       `json:"company_id" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uint      `json:"assignee_id"`
	Visibility  string     `json:"visibility" validate:"required,oneof=company private"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress done"`
}

func NewTaskFromCreateTaskByAdminRequestBody(createUserId uint, requestBody *CreateTaskByAdminRequestBody) *model.Task {
	return &model.Task{
		CompanyID:    requestBody.CompanyId,
		CreateUserId: createUserId,
		Title:        requestBody.Title,
		Description:  requestBody.Description,
		DueDate:      requestBody.DueDate,
		AssigneeID:   requestBody.AssigneeID,
		Visibility:   requestBody.Visibility,
		Status:       requestBody.Status,
	}
}

type UpdateTaskByAdminRequestBody struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uint      `json:"assignee_id"`
	Visibility  string     `json:"visibility" validate:"required,oneof=company private"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress done"`
}

func NewTaskFromUpdateTaskByAdminRequestBody(id uint, requestBody *UpdateTaskByAdminRequestBody) *model.Task {
	return &model.Task{
		ID:          id,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		DueDate:     requestBody.DueDate,
		AssigneeID:  requestBody.AssigneeID,
		Visibility:  requestBody.Visibility,
		Status:      requestBody.Status,
	}
}

type UpdateTaskRequestBody struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uint      `json:"assignee_id"`
	Visibility  string     `json:"visibility" validate:"required,oneof=company private"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress done"`
}

func NewTaskFromUpdateTaskRequestBody(id, companyId uint, requestBody *UpdateTaskRequestBody) *model.Task {
	return &model.Task{
		ID:          id,
		CompanyID:   companyId,
		Title:       requestBody.Title,
		Description: requestBody.Description,
		DueDate:     requestBody.DueDate,
		AssigneeID:  requestBody.AssigneeID,
		Visibility:  requestBody.Visibility,
		Status:      requestBody.Status,
	}
}
