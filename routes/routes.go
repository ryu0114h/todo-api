package routes

import (
	controller "todo-api/controller"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	apiV1 := e.Group("/api/v1")
	apiV1.GET("/tasks", controller.GetTasks)
	apiV1.GET("/tasks/:id", controller.GetTaskByID)
	apiV1.POST("/tasks", controller.CreateTask)
	apiV1.PATCH("/tasks/:id", controller.UpdateTask)
	apiV1.DELETE("/tasks/:id", controller.DeleteTask)
}
