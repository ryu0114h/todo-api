package routes

import (
	controller "todo-api/controller"
	"todo-api/repository"
	"todo-api/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	taskRepository := repository.NewTaskRepository(db)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	taskController := controller.NewTaskController(taskUseCase)

	apiV1 := e.Group("/api/v1")
	apiV1.GET("/tasks", taskController.GetTasks)
	apiV1.GET("/tasks/:id", taskController.GetTask)
	apiV1.POST("/tasks", taskController.CreateTask)
	apiV1.PATCH("/tasks/:id", taskController.UpdateTask)
	apiV1.DELETE("/tasks/:id", taskController.DeleteTask)
}
