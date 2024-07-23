package routes

import (
	controller "todo-api/controller"
	"todo-api/middleware"
	"todo-api/repository"
	"todo-api/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	taskRepository := repository.NewTaskRepository(db)
	companyRepository := repository.NewCompanyRepository(db)
	companyUserRepository := repository.NewCompanyUserRepository(db)
	taskUseCase := usecase.NewTaskUseCase(taskRepository, companyRepository, companyUserRepository)
	taskController := controller.NewTaskController(taskUseCase)

	apiV1 := e.Group("/api/v1")
	apiV1.Use(middleware.Logging())
	apiV1.GET("/companies/:company_id/tasks", taskController.GetTasks)
	apiV1.GET("/companies/:company_id/tasks/:task_id", taskController.GetTask)
	apiV1.POST("/companies/:company_id/tasks", taskController.CreateTask)
	apiV1.PUT("/companies/:company_id/tasks/:task_id", taskController.UpdateTask)
	apiV1.DELETE("/companies/:company_id/tasks/:task_id", taskController.DeleteTask)
}
