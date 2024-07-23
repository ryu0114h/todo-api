package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(companyId uint, limit, offset int) ([]*model.Task, error)
	GetTask(id, companyId uint) (*model.Task, error)
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(id uint, task *model.Task) (*model.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTasks(companyId uint, limit, offset int) ([]*model.Task, error) {
	tasks := []*model.Task{}
	result := r.db.Preload("Assignee").Where("company_id = ?", companyId).Limit(limit).Offset(offset).Find(&tasks)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTasks: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	return tasks, nil
}

func (r *taskRepository) GetTask(id, companyId uint) (*model.Task, error) {
	task := &model.Task{}
	result := r.db.Preload("Assignee").Where("company_id = ?", companyId).Find(task, "id = ?", id)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTask: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return task, nil
}

func (r *taskRepository) CreateTask(task *model.Task) (*model.Task, error) {
	if err := r.db.Create(task).Error; err != nil {
		return nil, myErrors.ErrDb
	}
	return task, nil
}

func (r *taskRepository) UpdateTask(id uint, task *model.Task) (*model.Task, error) {
	result := r.db.Save(task)
	if err := result.Error; err != nil {
		return nil, myErrors.ErrDb
	}
	return task, nil
}
