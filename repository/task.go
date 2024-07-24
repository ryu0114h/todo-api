package repository

import (
	"fmt"
	"log/slog"
	myErrors "todo-api/errors"
	"todo-api/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasksByCompanyId(companyId uint, limit, offset int) ([]*model.Task, error)
	GetTasks(limit, offset int) ([]*model.Task, error)
	GetTaskById(id uint) (*model.Task, error)
	GetTask(companyId, id uint) (*model.Task, error)
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(id uint, task *model.Task) (*model.Task, error)
	DeleteTaskById(id uint) error
	DeleteTask(companyId, id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetTasksByCompanyId(companyId uint, limit, offset int) ([]*model.Task, error) {
	tasks := []*model.Task{}
	result := r.db.Preload("Assignee").Where("company_id = ?", companyId).Limit(limit).Offset(offset).Find(&tasks)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTasks: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	return tasks, nil
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

func (r *taskRepository) GetTaskById(id uint) (*model.Task, error) {
	task := &model.Task{}
	result := r.db.Preload("Assignee").Where("id = ?", id).Find(task)
	if result.Error != nil {
		slog.Info(fmt.Sprintf("error GetTaskById: %v", result.Error))
		return nil, myErrors.ErrDb
	}
	if result.RowsAffected == 0 {
		return nil, myErrors.ErrNotFound
	}
	return task, nil
}

func (r *taskRepository) GetTask(companyId, id uint) (*model.Task, error) {
	task := &model.Task{}
	result := r.db.Preload("Assignee").Where("id = ? AND company_id = ?", id, companyId).Find(task)
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

func (r *taskRepository) DeleteTaskById(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return myErrors.ErrNotFound
	}
	return nil
}

func (r *taskRepository) DeleteTask(companyId, id uint) error {
	result := r.db.Where("id = ? AND company_id = ?", id, companyId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return myErrors.ErrNotFound
	}
	return nil
}
