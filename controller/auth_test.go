package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo-api/controller/request"
	"todo-api/controller/response"
	mock_usecase "todo-api/mock/usecase"
	"todo-api/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockAuthUseCase(ctrl)
	validate := validator.New()
	authController := NewAuthController(validate, mockUseCase)

	e := echo.New()

	tests := []struct {
		name           string
		requestBody    request.LoginRequestBody
		mockFunc       func()
		expectedStatus int
		expectedBody   *response.LoginResponseBody
	}{
		{
			name: "Success",
			requestBody: request.LoginRequestBody{
				Username: "testuser",
				Password: "password",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().Login("testuser", "password").Return(&model.User{ID: 1}, nil)
				os.Setenv("JWT_SECRET", "secret")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   &response.LoginResponseBody{Token: "<token>"},
		},
		{
			name: "Validation Error",
			requestBody: request.LoginRequestBody{
				Username: "",
				Password: "password",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "Invalid Request",
			requestBody: request.LoginRequestBody{
				Username: "",
				Password: "",
			},
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "Auth UseCase Error",
			requestBody: request.LoginRequestBody{
				Username: "testuser",
				Password: "password",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().Login("testuser", "password").Return(nil, errors.New("some error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "JWT Secret Not Found",
			requestBody: request.LoginRequestBody{
				Username: "testuser",
				Password: "password",
			},
			mockFunc: func() {
				mockUseCase.EXPECT().Login("testuser", "password").Return(&model.User{ID: 1}, nil)
				os.Setenv("JWT_SECRET", "")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// Assert the response
			if assert.NoError(t, authController.Login(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				if tt.expectedBody != nil {
					var actualBody response.LoginResponseBody
					if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actualBody)) {
						assert.NotNil(t, tt.expectedBody.Token, actualBody.Token)
					}
				}
			}
		})
	}
}
