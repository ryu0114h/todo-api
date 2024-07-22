package usecase

import (
	"todo-api/model"
	"todo-api/repository"
)

type TaskUseCase interface {
	GetTasks(limit, offset int) ([]*model.Task, error)
	GetTaskByID(id uint) (*model.Task, error)
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(id uint, task *model.Task) (*model.Task, error)
	DeleteTask(id uint) error
}

type taskUseCase struct {
	taskRepository repository.TaskRepository
}

func NewTaskUseCase(taskRepository repository.TaskRepository) TaskUseCase {
	return &taskUseCase{
		taskRepository: taskRepository,
	}
}

func (u *taskUseCase) GetTasks(limit, offset int) ([]*model.Task, error) {
	tasks, err := u.taskRepository.GetTasks(limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *taskUseCase) GetTaskByID(id uint) (*model.Task, error) {
	return nil, nil
}

func (u *taskUseCase) CreateTask(task *model.Task) (*model.Task, error) {
	return nil, nil
}

func (u *taskUseCase) UpdateTask(id uint, task *model.Task) (*model.Task, error) {
	return nil, nil
}

func (u *taskUseCase) DeleteTask(id uint) error {
	return nil
}
