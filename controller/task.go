package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetTasks(c echo.Context) error {
	return c.JSON(http.StatusOK, "Tasks")
}

func GetTaskByID(c echo.Context) error {
	return c.JSON(http.StatusOK, "Task by ID")
}

func CreateTask(c echo.Context) error {
	return c.JSON(http.StatusCreated, "Task Created")
}

func UpdateTask(c echo.Context) error {
	return c.JSON(http.StatusOK, "Task Updated")
}

func DeleteTask(c echo.Context) error {
	return c.JSON(http.StatusOK, "Task Deleted")
}
