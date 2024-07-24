package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/controller"
	"todo-api/controller/request"
	"todo-api/controller/response"
	mock_usecase "todo-api/mock/usecase"
	"todo-api/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserController_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockUserUseCase(ctrl)
	validate := validator.New()
	userController := controller.NewUserController(validate, mockUseCase)

	e := echo.New()

	tests := []struct {
		name           string
		requestBody    request.CreateUserRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Success",
			requestBody: request.CreateUserRequestBody{
				Username:   "testuser",
				Email:      "test@example.com",
				Password:   "password",
				Role:       "user",
				CompanyIds: []uint{1, 2},
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateUser(
					"testuser",
					"test@example.com",
					"password",
					"user",
					[]uint{1, 2},
				).Return(&model.User{ID: 1, Username: "testuser"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: &response.CreateUserResponseBody{
				User: &response.CreateUserResponseBodyUser{
					ID:       1,
					Username: "testuser",
				},
			},
		},
		{
			name: "Validation Error",
			requestBody: request.CreateUserRequestBody{
				Username: "",
				Email:    "test@example.com",
				Password: "password",
				Role:     "user",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "Key: 'CreateUserRequestBody.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'CreateUserRequestBody.CompanyIds' Error:Field validation for 'CompanyIds' failed on the 'required' tag"},
		},
		{
			name: "UseCase Error",
			requestBody: request.CreateUserRequestBody{
				Username:   "testuser",
				Email:      "test@example.com",
				Password:   "password",
				Role:       "user",
				CompanyIds: []uint{1, 2},
			},
			mockFunc: func() {
				mockUseCase.EXPECT().CreateUser(
					"testuser",
					"test@example.com",
					"password",
					"user",
					[]uint{1, 2},
				).Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "could not create user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			if assert.NoError(t, userController.CreateUser(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.expectedBody != nil {
					expectedJSON, _ := json.Marshal(tt.expectedBody)
					assert.JSONEq(t, string(expectedJSON), rec.Body.String())
				}
			}
		})
	}
}
