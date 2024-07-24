package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/controller/request"
	"todo-api/controller/response"
	myErrors "todo-api/errors"
	mock_usecase "todo-api/mock/usecase"
	"todo-api/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskController_GetTasksByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		limit          string
		offset         string
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			limit:  "10",
			offset: "0",
			mockFunc: func() {
				mockTasks := []*model.Task{
					{
						ID:          1,
						CompanyID:   1,
						Title:       "Task 1",
						Description: "Task 1 description",
						DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						AssigneeID:  &[]uint{1}[0],
						Visibility:  "company",
						Status:      "pending",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee: &model.User{
							ID:           11,
							Username:     "user 11",
							Email:        "user11@example.com",
							PasswordHash: "password",
							Role:         "user",
							CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
							UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						},
					},
					{
						ID:          2,
						CompanyID:   1,
						Title:       "Task 2",
						Description: "Task 2 description",
						DueDate:     nil,
						AssigneeID:  nil,
						Visibility:  "private",
						Status:      "in_progress",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee:    nil,
					},
				}
				mockUseCase.EXPECT().GetTasks(10, 0).Return(mockTasks, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: response.NewGetTasksResponseBody(
				[]*model.Task{
					{
						ID:          1,
						CompanyID:   1,
						Title:       "Task 1",
						Description: "Task 1 description",
						DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						AssigneeID:  &[]uint{1}[0],
						Visibility:  "company",
						Status:      "pending",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee: &model.User{
							ID:           11,
							Username:     "user 11",
							Email:        "user11@example.com",
							PasswordHash: "password",
							Role:         "user",
							CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
							UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						},
					},
					{
						ID:          2,
						CompanyID:   1,
						Title:       "Task 2",
						Description: "Task 2 description",
						DueDate:     nil,
						AssigneeID:  nil,
						Visibility:  "private",
						Status:      "in_progress",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee:    nil,
					},
				},
			),
		},
		{
			name:           "Limit exceeds max",
			limit:          "30",
			offset:         "0",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Limit exceeds the maximum allowed value of 20"},
		},
		{
			name:   "NotFound",
			limit:  "10",
			offset: "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasks(10, 0).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:   "InternalServerError",
			limit:  "10",
			offset: "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasks(10, 0).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/tasks?limit="+tc.limit+"&offset="+tc.offset, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			tc.mockFunc()

			if assert.NoError(t, taskController.GetTasksByAdmin(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_GetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyId      string
		limit          string
		offset         string
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyId: "1",
			limit:     "10",
			offset:    "0",
			mockFunc: func() {
				mockTasks := []*model.Task{
					{
						ID:          1,
						CompanyID:   1,
						Title:       "Task 1",
						Description: "Task 1 description",
						DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						AssigneeID:  &[]uint{1}[0],
						Visibility:  "company",
						Status:      "pending",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee: &model.User{
							ID:           11,
							Username:     "user 11",
							Email:        "user11@example.com",
							PasswordHash: "password",
							Role:         "user",
							CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
							UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						},
					},
					{
						ID:          2,
						CompanyID:   1,
						Title:       "Task 2",
						Description: "Task 2 description",
						DueDate:     nil,
						AssigneeID:  nil,
						Visibility:  "private",
						Status:      "in_progress",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee:    nil,
					},
				}
				mockUseCase.EXPECT().GetTasksByCompanyId(uint(1), 10, 0).Return(mockTasks, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: response.NewGetTasksResponseBody(
				[]*model.Task{
					{
						ID:          1,
						CompanyID:   1,
						Title:       "Task 1",
						Description: "Task 1 description",
						DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						AssigneeID:  &[]uint{1}[0],
						Visibility:  "company",
						Status:      "pending",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee: &model.User{
							ID:           11,
							Username:     "user 11",
							Email:        "user11@example.com",
							PasswordHash: "password",
							Role:         "user",
							CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
							UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						},
					},
					{
						ID:          2,
						CompanyID:   1,
						Title:       "Task 2",
						Description: "Task 2 description",
						DueDate:     nil,
						AssigneeID:  nil,
						Visibility:  "private",
						Status:      "in_progress",
						CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:   &[]time.Time{time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						Assignee:    nil,
					},
				},
			),
		},
		{
			name:           "BadRequest - Invalid CompanyID",
			companyId:      "invalid",
			limit:          "10",
			offset:         "0",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "company_id is bad request"},
		},
		{
			name:           "Limit exceeds max",
			companyId:      "1",
			limit:          "30",
			offset:         "0",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Limit exceeds the maximum allowed value of 20"},
		},
		{
			name:      "NotFound",
			companyId: "1",
			limit:     "10",
			offset:    "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasksByCompanyId(uint(1), 10, 0).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:      "InternalServerError",
			companyId: "1",
			limit:     "10",
			offset:    "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasksByCompanyId(uint(1), 10, 0).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+tc.companyId+"/tasks?limit="+tc.limit+"&offset="+tc.offset, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id")
			ctx.SetParamValues(tc.companyId)

			tc.mockFunc()

			if assert.NoError(t, taskController.GetTasks(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyID      string
		taskID         string
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				task := &model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Task 1",
					Description: "Task 1 description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
					CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					Assignee: &model.User{
						ID:           11,
						Username:     "user 11",
						Email:        "user11@example.com",
						PasswordHash: "password",
						Role:         "user",
						CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					},
				}
				mockUseCase.EXPECT().GetTask(uint(1), uint(2)).Return(task, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &response.GetTaskResponseBody{
				Task: &response.GetTaskResponseBodyTask{
					ID:          1,
					CompanyID:   1,
					Title:       "Task 1",
					Description: "Task 1 description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					Visibility:  "company",
					Status:      "pending",
					CreatedAt:   &[]time.Time{time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					UpdatedAt:   &[]time.Time{time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					Assignee: &response.GetTaskResponseBodyAssignee{
						ID:           11,
						Username:     "user 11",
						Email:        "user11@example.com",
						PasswordHash: "password",
						Role:         "user",
						CreatedAt:    &[]time.Time{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
						UpdatedAt:    &[]time.Time{time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					},
				},
			},
		},
		{
			name:      "Task not found",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTask(uint(1), uint(2)).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]string{
				"error": "not found",
			},
		},
		{
			name:           "Invalid task ID",
			companyID:      "1",
			taskID:         "invalid",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:           "Invalid company ID",
			companyID:      "invalid",
			taskID:         "1",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:      "Internal server error",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTask(uint(1), uint(2)).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "Failed to get task",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+tc.companyID+"/tasks/"+tc.taskID, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id", "task_id")
			ctx.SetParamValues(tc.companyID, tc.taskID)

			tc.mockFunc()

			if assert.NoError(t, taskController.GetTask(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)

				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_CreateTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		requestBody    *request.CreateTaskByAdminRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			requestBody: &request.CreateTaskByAdminRequestBody{
				CompanyId:   1,
				Title:       "New Task",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateTaskByAdmin(&model.Task{
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(&model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}, nil).Times(1)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &response.CreateTaskResponseBody{
				Task: &response.CreateTaskResponseBodyTask{
					ID:          1,
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				},
			},
		},
		{
			name: "Validation Error",
			requestBody: &request.CreateTaskByAdminRequestBody{
				CompanyId:   1,
				Title:       "",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Key: 'CreateTaskByAdminRequestBody.Title' Error:Field validation for 'Title' failed on the 'required' tag"},
		},
		{
			name: "Internal server error",
			requestBody: &request.CreateTaskByAdminRequestBody{
				CompanyId:   1,
				Title:       "New Task",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateTaskByAdmin(&model.Task{
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Failed to create task"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			tc.mockFunc()

			if assert.NoError(t, taskController.CreateTaskByAdmin(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyID      string
		requestBody    *request.CreateTaskRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyID: "1",
			requestBody: &request.CreateTaskRequestBody{
				Title:       "New Task",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateTask(&model.Task{
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(&model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}, nil).Times(1)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &response.CreateTaskResponseBody{
				Task: &response.CreateTaskResponseBodyTask{
					ID:          1,
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				},
			},
		},
		{
			name:      "Invalid company ID",
			companyID: "invalid",
			requestBody: &request.CreateTaskRequestBody{
				Title:       "New Task",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:      "Validation Error",
			companyID: "1",
			requestBody: &request.CreateTaskRequestBody{
				Title:       "",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Key: 'CreateTaskRequestBody.Title' Error:Field validation for 'Title' failed on the 'required' tag"},
		},
		{
			name:      "Internal server error",
			companyID: "1",
			requestBody: &request.CreateTaskRequestBody{
				Title:       "New Task",
				Description: "New task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateTask(&model.Task{
					CompanyID:   1,
					Title:       "New Task",
					Description: "New task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Failed to create task"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/companies/"+tc.companyID+"/tasks", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id")
			ctx.SetParamValues(tc.companyID)

			tc.mockFunc()

			if assert.NoError(t, taskController.CreateTask(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_UpdateTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		taskID         string
		requestBody    *request.UpdateTaskByAdminRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Success",
			taskID: "1",
			requestBody: &request.UpdateTaskByAdminRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().UpdateTaskByAdmin(uint(1), &model.Task{
					ID:          1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(&model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &response.UpdateTaskResponseBody{
				Task: &response.UpdateTaskResponseBodyTask{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				},
			},
		},
		{
			name:   "Invalid task ID",
			taskID: "invalid",
			requestBody: &request.UpdateTaskByAdminRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:   "Validation Error",
			taskID: "1",
			requestBody: &request.UpdateTaskByAdminRequestBody{
				Title:       "",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Key: 'UpdateTaskByAdminRequestBody.Title' Error:Field validation for 'Title' failed on the 'required' tag"},
		},
		{
			name:   "Not found",
			taskID: "1",
			requestBody: &request.UpdateTaskByAdminRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {
				mockUseCase.EXPECT().UpdateTaskByAdmin(uint(1), &model.Task{
					ID:          1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:   "Internal server error",
			taskID: "1",
			requestBody: &request.UpdateTaskByAdminRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {
				mockUseCase.EXPECT().UpdateTaskByAdmin(uint(1), &model.Task{
					ID:          1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/tasks/"+tc.taskID, bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("task_id")
			ctx.SetParamValues(tc.taskID)

			tc.mockFunc()

			if assert.NoError(t, taskController.UpdateTaskByAdmin(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyID      string
		taskID         string
		requestBody    *request.UpdateTaskRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyID: "1",
			taskID:    "1",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().UpdateTask(uint(1), uint(1), &model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(&model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &response.UpdateTaskResponseBody{
				Task: &response.UpdateTaskResponseBodyTask{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				},
			},
		},
		{
			name:      "Invalid task ID",
			companyID: "1",
			taskID:    "invalid",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:      "Invalid company ID",
			companyID: "invalid",
			taskID:    "1",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:      "Validation Error",
			companyID: "1",
			taskID:    "1",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Key: 'UpdateTaskRequestBody.Title' Error:Field validation for 'Title' failed on the 'required' tag"},
		},
		{
			name:      "Not found",
			companyID: "1",
			taskID:    "1",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {
				mockUseCase.EXPECT().UpdateTask(uint(1), uint(1), &model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:      "Internal server error",
			companyID: "1",
			taskID:    "1",
			requestBody: &request.UpdateTaskRequestBody{
				Title:       "Updated Task",
				Description: "Updated task description",
				DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
				AssigneeID:  &[]uint{1}[0],
				Visibility:  "company",
				Status:      "pending",
			}, mockFunc: func() {
				mockUseCase.EXPECT().UpdateTask(uint(1), uint(1), &model.Task{
					ID:          1,
					CompanyID:   1,
					Title:       "Updated Task",
					Description: "Updated task description",
					DueDate:     &[]time.Time{time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)}[0],
					AssigneeID:  &[]uint{1}[0],
					Visibility:  "company",
					Status:      "pending",
				}).Return(nil, errors.New("some error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPut, "/api/v1/companies/"+tc.companyID+"/tasks/"+tc.taskID, bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id", "task_id")
			ctx.SetParamValues(tc.companyID, tc.taskID)

			tc.mockFunc()

			if assert.NoError(t, taskController.UpdateTask(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_DeleteTaskByAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyID      string
		taskID         string
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTaskByAdmin(uint(2)).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Invalid task ID",
			companyID:      "1",
			taskID:         "invalid",
			mockFunc:       func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name:      "Not found",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTaskByAdmin(uint(2)).Return(myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:      "Internal server error",
			companyID: "1",
			taskID:    "1",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTaskByAdmin(uint(1)).Return(errors.New("internal error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/tasks/"+tc.taskID, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id", "task_id")
			ctx.SetParamValues(tc.companyID, tc.taskID)

			tc.mockFunc()

			if assert.NoError(t, taskController.DeleteTaskByAdmin(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}

func TestTaskController_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	validate := validator.New()
	taskController := NewTaskController(validate, mockUseCase)

	testCases := []struct {
		name           string
		companyID      string
		taskID         string
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "Success",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTask(uint(1), uint(2)).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Invalid task ID",
			companyID:      "1",
			taskID:         "invalid",
			mockFunc:       func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name:           "Invalid company ID",
			companyID:      "invalid",
			taskID:         "2",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:      "Not found",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTask(uint(1), uint(2)).Return(myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:      "Internal server error",
			companyID: "1",
			taskID:    "2",
			mockFunc: func() {
				mockUseCase.EXPECT().DeleteTask(uint(1), uint(2)).Return(errors.New("internal error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/companies/"+tc.companyID+"/tasks/"+tc.taskID, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("company_id", "task_id")
			ctx.SetParamValues(tc.companyID, tc.taskID)

			tc.mockFunc()

			if assert.NoError(t, taskController.DeleteTask(ctx)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tc.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}
