package response

import (
	"time"
	"todo-api/model"
)

// GetTasksResponseBody はタスク一覧取得リクエストのレスポンスボディ
type GetTasksResponseBody struct {
	Tasks []*GetTasksResponseBodyTask `json:"tasks"`
}

type GetTasksResponseBodyTask struct {
	ID          uint                          `json:"id"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	DueDate     *time.Time                    `json:"due_date"`
	Visibility  string                        `json:"visibility"`
	Status      string                        `json:"status"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	Assignee    *GetTasksResponseBodyAssignee `json:"assignee"`
}

type GetTasksResponseBodyAssignee struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewGetTasksResponseBody(tasks []*model.Task) *GetTasksResponseBody {
	resTasks := []*GetTasksResponseBodyTask{}

	for _, task := range tasks {
		resTasks = append(resTasks, &GetTasksResponseBodyTask{
			ID:          task.ID,
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
