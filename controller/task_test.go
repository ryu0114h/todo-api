package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo-api/controller/response"
	myErrors "todo-api/errors"
	mock_usecase "todo-api/mock/usecase"
	"todo-api/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskController_GetTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTaskUseCase(ctrl)
	taskController := NewTaskController(mockUseCase)

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
				mockUseCase.EXPECT().GetTasks(uint(1), 10, 0).Return(mockTasks, nil).Times(1)
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
			name:      "NotFound",
			companyId: "1",
			limit:     "10",
			offset:    "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasks(uint(1), 10, 0).Return(nil, myErrors.ErrNotFound).Times(1)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "not found"},
		},
		{
			name:      "InternalError",
			companyId: "1",
			limit:     "10",
			offset:    "0",
			mockFunc: func() {
				mockUseCase.EXPECT().GetTasks(uint(1), 10, 0).Return(nil, errors.New("some error")).Times(1)
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
	taskController := NewTaskController(mockUseCase)

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
