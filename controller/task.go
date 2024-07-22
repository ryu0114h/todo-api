package controller

import (
	"net/http"
	"strconv"
	"todo-api/model"
	"todo-api/usecase"

	"github.com/labstack/echo/v4"
)

const (
	// タスクのデフォルトの取得制限数
	DEFAULT_TASK_LIMIT = 10

	// タスクのデフォルトのオフセット値
	DEFAULT_TASK_OFFSET = 0
)

type TaskController interface {
	GetTasks(ctx echo.Context) error
	GetTaskByID(ctx echo.Context) error
	CreateTask(ctx echo.Context) error
	UpdateTask(ctx echo.Context) error
	DeleteTask(ctx echo.Context) error
}

type taskController struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskController(taskUseCase usecase.TaskUseCase) TaskController {
	return &taskController{taskUseCase: taskUseCase}
}

func (c *taskController) GetTasks(ctx echo.Context) error {
	limitStr := ctx.QueryParam("limit")
	offsetStr := ctx.QueryParam("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = DEFAULT_TASK_LIMIT
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = DEFAULT_TASK_OFFSET
	}

	tasks, err := c.taskUseCase.GetTasks(limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}
	return ctx.JSON(http.StatusOK, tasks)
}

func (c *taskController) GetTaskByID(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	task, err := c.taskUseCase.GetTaskByID(uint(id))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}
	return ctx.JSON(http.StatusOK, task)
}

func (c *taskController) CreateTask(ctx echo.Context) error {
	task, err := c.taskUseCase.CreateTask(&model.Task{})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusCreated, task)
}

func (c *taskController) UpdateTask(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	task, err := c.taskUseCase.UpdateTask(uint(id), &model.Task{})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, task)
}

func (c *taskController) DeleteTask(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	err = c.taskUseCase.DeleteTask(uint(id))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, nil)
}
