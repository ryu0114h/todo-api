package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(limit, offset int) ([]*model.Task, error)
	GetTask(id uint) (*model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTasks(limit, offset int) ([]*model.Task, error) {
	tasks := []*model.Task{}
	result := r.db.Preload("Assignee").Limit(limit).Offset(offset).Find(&tasks)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTasks: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	return tasks, nil
}

func (r *taskRepository) GetTask(id uint) (*model.Task, error) {
	task := &model.Task{}
	result := r.db.Preload("Assignee").Find(task, "id = ?", id)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTask: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return task, nil
}
