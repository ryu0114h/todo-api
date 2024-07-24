package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	mock_repository "todo-api/mock/repository"
	"todo-api/model"
)

func TestAuthUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	authUseCase := NewAuthUseCase(mockUserRepo)

	// テスト用ユーザー
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	testCases := []struct {
		name                  string
		username              string
		password              string
		mockGetUserByUsername func()
		expectedResult        *model.User
		expectedError         error
	}{
		{
			name:     "Success",
			username: "testuser",
			password: "password",
			mockGetUserByUsername: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("testuser").
					Return(user, nil).
					Times(1)
			},
			expectedResult: user,
			expectedError:  nil,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			password: "password",
			mockGetUserByUsername: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("nonexistent").
					Return(nil, errors.New("user not found")).
					Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:     "Error from GetUserByUsername",
			username: "testuser",
			password: "password",
			mockGetUserByUsername: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("testuser").
					Return(nil, errors.New("db error")).
					Times(1)
			},
			expectedResult: nil,
			expectedError:  errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockGetUserByUsername()

			// `Login` メソッドの呼び出し
			result, err := authUseCase.Login(tc.username, tc.password)

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
