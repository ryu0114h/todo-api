package routes

import (
	controller "todo-api/controller"
	"todo-api/middleware"
	"todo-api/repository"
	"todo-api/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	var validate = validator.New()
	taskRepository := repository.NewTaskRepository(db)
	companyRepository := repository.NewCompanyRepository(db)
	companyUserRepository := repository.NewCompanyUserRepository(db)
	userRepository := repository.NewUserRepository(db)
	taskUseCase := usecase.NewTaskUseCase(taskRepository, companyRepository, companyUserRepository)
	userUseCase := usecase.NewUserUseCase(db, userRepository, companyUserRepository)
	authUseCase := usecase.NewAuthUseCase(userRepository)
	taskController := controller.NewTaskController(validate, taskUseCase)
	userController := controller.NewUserController(validate, userUseCase)
	authController := controller.NewAuthController(validate, authUseCase)

	apiV1 := e.Group("/api/v1")
	apiV1.Use(middleware.Logging())
	apiV1.POST("/users", userController.CreateUser)
	apiV1.POST("/login", authController.Login)

	// 下記は認証が必要なAPI
	apiV1.Use(middleware.Auth(userRepository))
	apiV1Company := apiV1.Group("/companies/:company_id")
	apiV1Company.Use(middleware.CompanyAuth(companyUserRepository))
	apiV1Company.GET("/tasks", taskController.GetTasks)
	apiV1Company.GET("/tasks/:task_id", taskController.GetTask)
	apiV1Company.POST("/tasks", taskController.CreateTask)
	apiV1Company.PUT("/tasks/:task_id", taskController.UpdateTask)
	apiV1Company.DELETE("/tasks/:task_id", taskController.DeleteTask)

	// 下記はAdminユーザ専用のAPI
	apiV1Admin := apiV1.Group("/admin")
	apiV1Admin.Use(middleware.AdminAuth())
	apiV1Admin.GET("/tasks", taskController.GetTasksByAdmin)
	apiV1Admin.POST("/tasks", taskController.CreateTaskByAdmin)
	apiV1Admin.PUT("/tasks/:task_id", taskController.UpdateTaskByAdmin)
	apiV1Admin.DELETE("/tasks/:task_id", taskController.DeleteTaskByAdmin)
}
