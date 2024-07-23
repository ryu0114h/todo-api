package response

import (
	"time"
	"todo-api/model"
)

// GetTasksResponseBody はタスク一覧取得APIのレスポンスボディ
type GetTasksResponseBody struct {
	Tasks []*GetTasksResponseBodyTask `json:"tasks"`
}

type GetTasksResponseBodyTask struct {
	ID          uint                          `json:"id"`
	CompanyID   uint                          `json:"company_id"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	DueDate     *time.Time                    `json:"due_date"`
	Visibility  string                        `json:"visibility"`
	Status      string                        `json:"status"`
	CreatedAt   *time.Time                    `json:"created_at"`
	UpdatedAt   *time.Time                    `json:"updated_at"`
	Assignee    *GetTasksResponseBodyAssignee `json:"assignee"`
}

type GetTasksResponseBodyAssignee struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         string     `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func NewGetTasksResponseBody(tasks []*model.Task) *GetTasksResponseBody {
	resTasks := []*GetTasksResponseBodyTask{}

	for _, task := range tasks {
		resTasks = append(resTasks, &GetTasksResponseBodyTask{
			ID:          task.ID,
			CompanyID:   task.CompanyID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Visibility:  task.Visibility,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Assignee:    NewGetTasksResponseBodyAssignee(task.Assignee),
		})
	}

	return &GetTasksResponseBody{
		Tasks: resTasks,
	}
}

func NewGetTasksResponseBodyAssignee(assignee *model.User) *GetTasksResponseBodyAssignee {
	if assignee == nil {
		return nil
	}

	return &GetTasksResponseBodyAssignee{
		ID:           assignee.ID,
		Username:     assignee.Username,
		Email:        assignee.Email,
		PasswordHash: assignee.PasswordHash,
		Role:         assignee.Role,
		CreatedAt:    assignee.CreatedAt,
		UpdatedAt:    assignee.UpdatedAt,
	}
}

// GetTaskResponseBody はタスク取得APIのレスポンスボディ
type GetTaskResponseBody struct {
	Task *GetTaskResponseBodyTask `json:"task"`
}

type GetTaskResponseBodyTask struct {
	ID          uint                         `json:"id"`
	CompanyID   uint                         `json:"company_id"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	DueDate     *time.Time                   `json:"due_date"`
	Visibility  string                       `json:"visibility"`
	Status      string                       `json:"status"`
	CreatedAt   *time.Time                   `json:"created_at"`
	UpdatedAt   *time.Time                   `json:"updated_at"`
	Assignee    *GetTaskResponseBodyAssignee `json:"assignee"`
}

type GetTaskResponseBodyAssignee struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	Role         string     `json:"role"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func NewGetTaskResponseBody(task *model.Task) *GetTaskResponseBody {
	return &GetTaskResponseBody{
		Task: &GetTaskResponseBodyTask{
			ID:          task.ID,
			CompanyID:   task.CompanyID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Visibility:  task.Visibility,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Assignee:    NewGetTaskResponseBodyAssignee(task.Assignee),
		},
	}
}

func NewGetTaskResponseBodyAssignee(assignee *model.User) *GetTaskResponseBodyAssignee {
	if assignee == nil {
		return nil
	}

	return &GetTaskResponseBodyAssignee{
		ID:           assignee.ID,
		Username:     assignee.Username,
		Email:        assignee.Email,
		PasswordHash: assignee.PasswordHash,
		Role:         assignee.Role,
		CreatedAt:    assignee.CreatedAt,
		UpdatedAt:    assignee.UpdatedAt,
	}
}

// CreateTaskResponseBody はタスク取得APIのレスポンスボディ
type CreateTaskResponseBody struct {
	Task *CreateTaskResponseBodyTask `json:"task"`
}

type CreateTaskResponseBodyTask struct {
	ID          uint       `json:"id"`
	CompanyID   uint       `json:"company_id"`
	AssigneeID  *uint      `json:"assignee_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Visibility  string     `json:"visibility"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func NewCreateTaskResponseBody(task *model.Task) *CreateTaskResponseBody {
	return &CreateTaskResponseBody{
		Task: &CreateTaskResponseBodyTask{
			ID:          task.ID,
			CompanyID:   task.CompanyID,
			AssigneeID:  task.AssigneeID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Visibility:  task.Visibility,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	}
}

type UpdateTaskResponseBody struct {
	Task *UpdateTaskResponseBodyTask `json:"task"`
}

type UpdateTaskResponseBodyTask struct {
	ID          uint       `json:"id"`
	CompanyID   uint       `json:"company_id"`
	AssigneeID  *uint      `json:"assignee_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Visibility  string     `json:"visibility"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func NewUpdateTaskResponseBody(task *model.Task) *UpdateTaskResponseBody {
	return &UpdateTaskResponseBody{
		Task: &UpdateTaskResponseBodyTask{
			ID:          task.ID,
			CompanyID:   task.CompanyID,
			AssigneeID:  task.AssigneeID,
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Visibility:  task.Visibility,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		},
	}
}
