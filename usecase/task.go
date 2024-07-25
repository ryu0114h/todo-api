package usecase

import (
	"todo-api/model"
	"todo-api/repository"
)

type TaskUseCase interface {
	GetTasksByCompanyId(companyId, createUserId uint, limit, offset int) ([]*model.Task, error)
	GetTasks(limit, offset int) ([]*model.Task, error)
	GetTask(companyId, taskId, createUserId uint) (*model.Task, error)
	CreateTaskByAdmin(task *model.Task) (*model.Task, error)
	CreateTask(task *model.Task) (*model.Task, error)
	UpdateTaskByAdmin(taskId uint, task *model.Task) (*model.Task, error)
	UpdateTask(companyId, taskId, createUserId uint, task *model.Task) (*model.Task, error)
	DeleteTaskByAdmin(taskId uint) error
	DeleteTask(companyId, taskId, createUserId uint) error
}

type taskUseCase struct {
	taskRepository        repository.TaskRepository
	companyRepository     repository.CompanyRepository
	companyUserRepository repository.CompanyUserRepository
}

func NewTaskUseCase(
	taskRepository repository.TaskRepository,
	companyRepository repository.CompanyRepository,
	companyUserRepository repository.CompanyUserRepository,
) TaskUseCase {
	return &taskUseCase{
		taskRepository:        taskRepository,
		companyRepository:     companyRepository,
		companyUserRepository: companyUserRepository,
	}
}

func (u *taskUseCase) GetTasksByCompanyId(companyId, createUserId uint, limit, offset int) ([]*model.Task, error) {
	_, err := u.companyRepository.GetCompany(companyId)
	if err != nil {
		return nil, err
	}

	tasks, err := u.taskRepository.GetTasksByCompanyId(companyId, createUserId, limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *taskUseCase) GetTasks(limit, offset int) ([]*model.Task, error) {
	tasks, err := u.taskRepository.GetTasks(limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u *taskUseCase) GetTask(companyId, taskId, createUserId uint) (*model.Task, error) {
	task, err := u.taskRepository.GetTask(companyId, taskId, createUserId)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUseCase) CreateTaskByAdmin(task *model.Task) (*model.Task, error) {
	task, err := u.taskRepository.CreateTask(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUseCase) CreateTask(task *model.Task) (*model.Task, error) {
	if task.AssigneeID != nil {
		_, err := u.companyUserRepository.GetCompanyUser(task.CompanyID, *task.AssigneeID)
		if err != nil {
			return nil, err
		}
	}

	task, err := u.taskRepository.CreateTask(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUseCase) UpdateTaskByAdmin(taskId uint, task *model.Task) (*model.Task, error) {
	oldTask, err := u.taskRepository.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}

	resultTask, err := u.taskRepository.UpdateTask(taskId, &model.Task{
		ID:           oldTask.ID,
		CompanyID:    oldTask.CompanyID,
		CreateUserId: oldTask.CreateUserId,
		Title:        task.Title,
		Description:  task.Description,
		DueDate:      task.DueDate,
		AssigneeID:   task.AssigneeID,
		Visibility:   task.Visibility,
		Status:       task.Status,
		CreatedAt:    oldTask.CreatedAt,
	})
	if err != nil {
		return nil, err
	}
	return resultTask, nil
}

func (u *taskUseCase) UpdateTask(companyId, taskId, createUserId uint, task *model.Task) (*model.Task, error) {
	if task.AssigneeID != nil {
		_, err := u.companyUserRepository.GetCompanyUser(task.CompanyID, *task.AssigneeID)
		if err != nil {
			return nil, err
		}
	}

	oldTask, err := u.taskRepository.GetTask(companyId, taskId, createUserId)
	if err != nil {
		return nil, err
	}

	resultTask, err := u.taskRepository.UpdateTask(taskId, &model.Task{
		ID:           oldTask.ID,
		CompanyID:    oldTask.CompanyID,
		CreateUserId: oldTask.CreateUserId,
		Title:        task.Title,
		Description:  task.Description,
		DueDate:      task.DueDate,
		AssigneeID:   task.AssigneeID,
		Visibility:   task.Visibility,
		Status:       task.Status,
		CreatedAt:    oldTask.CreatedAt,
	})
	if err != nil {
		return nil, err
	}
	return resultTask, nil
}

func (u *taskUseCase) DeleteTaskByAdmin(taskId uint) error {
	_, err := u.taskRepository.GetTaskById(taskId)
	if err != nil {
		return err
	}

	err = u.taskRepository.DeleteTaskById(taskId)
	if err != nil {
		return err
	}

	return nil
}

func (u *taskUseCase) DeleteTask(companyId, taskId, createUserId uint) error {
	_, err := u.taskRepository.GetTask(companyId, taskId, createUserId)
	if err != nil {
		return err
	}

	err = u.taskRepository.DeleteTask(companyId, taskId)
	if err != nil {
		return err
	}

	return nil
}
