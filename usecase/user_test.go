package usecase

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"todo-api/mock"
	mock_repository "todo-api/mock/repository"
	"todo-api/model"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックの設定
	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockCompanyUserRepo := mock_repository.NewMockCompanyUserRepository(ctrl)

	user := &model.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	testCases := []struct {
		name           string
		mock           func(sqlMock sqlmock.Sqlmock)
		expectedResult *model.User
		expectedError  error
	}{
		{
			name: "Success",
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				mockUserRepo.EXPECT().
					CreateUser(gomock.Any()).
					Return(user, nil).
					Times(1)
				mockCompanyUserRepo.EXPECT().
					CreateCompanyUsers(gomock.Any()).
					Return([]*model.CompanyUser{}, nil).
					Times(1)
				sqlMock.ExpectCommit()
			},
			expectedResult: user,
			expectedError:  nil,
		},
		{
			name: "Error in CreateUser",
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				mockUserRepo.EXPECT().
					CreateUser(gomock.Any()).
					Return(nil, errors.New("some error")).
					Times(1)
				sqlMock.ExpectExec("INSERT INTO `users`").
					WithArgs(gomock.Any()).
					WillReturnError(errors.New("some error"))
				sqlMock.ExpectRollback()
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
		{
			name: "Error in CreateCompanyUsers",
			mock: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				mockUserRepo.EXPECT().
					CreateUser(gomock.Any()).
					Return(user, nil).
					Times(1)
				sqlMock.ExpectExec("INSERT INTO `users`").
					WithArgs(gomock.Any()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mockCompanyUserRepo.EXPECT().
					CreateCompanyUsers(gomock.Any()).
					Return(nil, errors.New("some error")).
					Times(1)
				sqlMock.ExpectExec("INSERT INTO `company_users`").
					WithArgs(gomock.Any()).
					WillReturnError(errors.New("some error"))
				sqlMock.ExpectRollback()
			},
			expectedResult: nil,
			expectedError:  errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// デフォルトの空の *gorm.DB インスタンス
			db, sqlMock, err := mock.NewDbMock()
			assert.NoError(t, err)

			tc.mock(sqlMock)

			userUseCase := NewUserUseCase(db, mockUserRepo, mockCompanyUserRepo)
			createdUser, err := userUseCase.CreateUser("testuser", "test@example.com", "password", "user", []uint{1})

			assert.Equal(t, tc.expectedResult, createdUser)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
