package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_repository "todo-api/mock/repository"
	"todo-api/model"
	"todo-api/usecase"
)

func TestTaskUseCase_GetTasksByCompanyId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	companyId := uint(1)
	userId := uint(2)
	limit := 10
	offset := 0

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult []*model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockCompanyRepo.EXPECT().GetCompany(companyId).Return(&model.Company{ID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().GetTasksByCompanyId(companyId, userId, limit, offset).Return([]*model.Task{}, nil).Times(1)
			},
			expectedResult: []*model.Task{},
			expectedError:  nil,
		},
		{
			name: "Company not found",
			mockFunc: func() {
				mockCompanyRepo.EXPECT().GetCompany(companyId).Return(nil, errors.New("company not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("company not found"),
		},
		{
			name: "Error in GetTasksByCompanyId",
			mockFunc: func() {
				mockCompanyRepo.EXPECT().GetCompany(companyId).Return(&model.Company{ID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().GetTasksByCompanyId(companyId, userId, limit, offset).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			tasks, err := taskUseCase.GetTasksByCompanyId(companyId, userId, limit, offset)

			assert.Equal(t, tc.expectedResult, tasks)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_GetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	limit := 10
	offset := 0

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult []*model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTasks(limit, offset).Return([]*model.Task{}, nil).Times(1)
			},
			expectedResult: []*model.Task{},
			expectedError:  nil,
		},
		{
			name: "Error in GetTasks",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTasks(limit, offset).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			tasks, err := taskUseCase.GetTasks(limit, offset)

			assert.Equal(t, tc.expectedResult, tasks)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	companyId := uint(1)
	userId := uint(2)
	taskId := uint(3)

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult *model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(&model.Task{ID: taskId, CompanyID: companyId}, nil).Times(1)
			},
			expectedResult: &model.Task{ID: taskId, CompanyID: companyId},
			expectedError:  nil,
		},
		{
			name: "Task not found",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(nil, errors.New("task not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("task not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			task, err := taskUseCase.GetTask(companyId, taskId, userId)

			assert.Equal(t, tc.expectedResult, task)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_CreateTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	task := &model.Task{
		ID:          1,
		CompanyID:   1,
		Title:       "Task Title",
		Description: "Task Description",
		DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
		AssigneeID:  &[]uint{1}[0],
		Visibility:  "public",
		Status:      "pending",
	}

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult *model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().CreateTask(task).Return(task, nil).Times(1)
			},
			expectedResult: task,
			expectedError:  nil,
		},

		{
			name: "Error in CreateTask",
			mockFunc: func() {
				mockTaskRepo.EXPECT().CreateTask(task).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			createdTask, err := taskUseCase.CreateTaskByAdmin(task)

			assert.Equal(t, tc.expectedResult, createdTask)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	task := &model.Task{
		ID:          1,
		CompanyID:   1,
		Title:       "Task Title",
		Description: "Task Description",
		DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
		AssigneeID:  &[]uint{1}[0],
		Visibility:  "public",
		Status:      "pending",
	}

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult *model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(&model.CompanyUser{ID: *task.AssigneeID}, nil).Times(1)
				mockTaskRepo.EXPECT().CreateTask(task).Return(task, nil).Times(1)
			},
			expectedResult: task,
			expectedError:  nil,
		},
		{
			name: "Assignee not found",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(nil, errors.New("assignee not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("assignee not found"),
		},
		{
			name: "Error in CreateTask",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(&model.CompanyUser{ID: *task.AssigneeID}, nil).Times(1)
				mockTaskRepo.EXPECT().CreateTask(task).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			createdTask, err := taskUseCase.CreateTask(task)

			assert.Equal(t, tc.expectedResult, createdTask)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_UpdateTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	companyId := uint(1)
	taskId := uint(2)
	task := &model.Task{
		CompanyID:   companyId,
		Title:       "Updated Task Title",
		Description: "Updated Task Description",
		DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
		AssigneeID:  &[]uint{1}[0],
		Visibility:  "public",
		Status:      "completed",
	}

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult *model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(&model.Task{ID: taskId, CompanyID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().UpdateTask(taskId, &model.Task{
					ID:          taskId,
					CompanyID:   companyId,
					Title:       "Updated Task Title",
					Description: "Updated Task Description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "public",
					Status:      "completed",
				}).Return(task, nil).Times(1)
			},
			expectedResult: task,
			expectedError:  nil,
		},
		{
			name: "Task not found",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(nil, errors.New("task not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("task not found"),
		},
		{
			name: "Error in UpdateTask",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(&model.Task{ID: taskId, CompanyID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().UpdateTask(taskId, &model.Task{
					ID:          taskId,
					CompanyID:   companyId,
					Title:       "Updated Task Title",
					Description: "Updated Task Description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "public",
					Status:      "completed",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			updatedTask, err := taskUseCase.UpdateTaskByAdmin(taskId, task)

			assert.Equal(t, tc.expectedResult, updatedTask)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	companyId := uint(1)
	userId := uint(2)
	taskId := uint(3)
	task := &model.Task{
		CompanyID:   companyId,
		Title:       "Updated Task Title",
		Description: "Updated Task Description",
		DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
		AssigneeID:  &[]uint{1}[0],
		Visibility:  "public",
		Status:      "completed",
	}

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedResult *model.Task
		expectedError  error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(&model.CompanyUser{ID: *task.AssigneeID}, nil).Times(1)
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(&model.Task{ID: taskId, CompanyID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().UpdateTask(taskId, &model.Task{
					ID:          taskId,
					CompanyID:   companyId,
					Title:       "Updated Task Title",
					Description: "Updated Task Description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "public",
					Status:      "completed",
				}).Return(task, nil).Times(1)
			},
			expectedResult: task,
			expectedError:  nil,
		},
		{
			name: "Assignee not found",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(nil, errors.New("assignee not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("assignee not found"),
		},
		{
			name: "Task not found",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(&model.CompanyUser{ID: *task.AssigneeID}, nil).Times(1)
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(nil, errors.New("task not found")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("task not found"),
		},
		{
			name: "Error in UpdateTask",
			mockFunc: func() {
				mockCompanyUserRepo.EXPECT().GetCompanyUser(task.CompanyID, *task.AssigneeID).Return(&model.CompanyUser{ID: *task.AssigneeID}, nil).Times(1)
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(&model.Task{ID: taskId, CompanyID: companyId}, nil).Times(1)
				mockTaskRepo.EXPECT().UpdateTask(taskId, &model.Task{
					ID:          taskId,
					CompanyID:   companyId,
					Title:       "Updated Task Title",
					Description: "Updated Task Description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "public",
					Status:      "completed",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			updatedTask, err := taskUseCase.UpdateTask(companyId, taskId, userId, task)

			assert.Equal(t, tc.expectedResult, updatedTask)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_DeleteTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	taskId := uint(1)

	testCases := []struct {
		name          string
		mockFunc      func()
		expectedError error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(&model.Task{
					ID: taskId,
				}, nil).Times(1)
				mockTaskRepo.EXPECT().DeleteTaskById(taskId).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Error in GetTaskById",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(nil,
					errors.New("some error")).Times(1)
			},
			expectedError: errors.New("some error"),
		},
		{
			name: "Error in DeleteTaskById",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTaskById(taskId).Return(&model.Task{
					ID: taskId,
				}, nil).Times(1)
				mockTaskRepo.EXPECT().DeleteTaskById(taskId).Return(errors.New("some error")).Times(1)
			},
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := taskUseCase.DeleteTaskByAdmin(taskId)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTaskUseCase_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock_repository.NewMockTaskRepository(ctrl)
	mockCompanyRepo := mock_repository.NewMockCompanyRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	taskUseCase := usecase.NewTaskUseCase(mockTaskRepo, mockCompanyRepo, mockCompanyUserRepo)

	companyId := uint(1)
	taskId := uint(2)
	userId := uint(3)

	testCases := []struct {
		name          string
		mockFunc      func()
		expectedError error
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(&model.Task{
					ID: taskId,
				}, nil).Times(1)
				mockTaskRepo.EXPECT().DeleteTask(companyId, taskId).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Error in GetTask",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(nil,
					errors.New("some error")).Times(1)
			},
			expectedError: errors.New("some error"),
		},
		{
			name: "Error in DeleteTask",
			mockFunc: func() {
				mockTaskRepo.EXPECT().GetTask(companyId, taskId, userId).Return(&model.Task{
					ID: taskId,
				}, nil).Times(1)
				mockTaskRepo.EXPECT().DeleteTask(companyId, taskId).Return(errors.New("some error")).Times(1)
			},
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := taskUseCase.DeleteTask(companyId, taskId, userId)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
