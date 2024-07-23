package usecase

import (
	"todo-api/model"
	"todo-api/repository"
)

type TaskUseCase interface {
	GetTasks(companyId uint, limit, offset int) ([]*model.Task, error)
	GetTask(id, companyId uint) (*model.Task, error)
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTask(id uint, task *model.Task) (*model.Task, error)
	DeleteTask(id uint) error
}

type taskUseCase struct {
	taskRepository    repository.TaskRepository
	companyRepository repository.CompanyRepository
}

func NewTaskUseCase(taskRepository repository.TaskRepository, companyRepository repository.CompanyRepository) TaskUseCase {
	return &taskUseCase{
		taskRepository:    taskRepository,
		companyRepository: companyRepository,
	}
}

func (u *taskUseCase) GetTasks(companyId uint, limit, offset int) ([]*model.Task, error) {
	_, err := u.companyRepository.GetCompany(companyId)
	if err != nil {
		return nil, err
	}

	tasks, err := u.taskRepository.GetTasks(companyId, limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *taskUseCase) GetTask(id, companyId uint) (*model.Task, error) {
	task, err := u.taskRepository.GetTask(companyId, id)
	if err != nil {
		return nil, err
	}
	return task, nil
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
