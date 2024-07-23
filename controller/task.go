package controller

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"todo-api/controller/request"
	"todo-api/controller/response"
	myErrors "todo-api/errors"
	"todo-api/usecase"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

const (
	// タスクのデフォルトの取得制限数
	DEFAULT_TASK_LIMIT = 10

	// タスクのデフォルトのオフセット値
	DEFAULT_TASK_OFFSET = 0
)

var validate = validator.New()

type TaskController interface {
	GetTasks(ctx echo.Context) error
	GetTask(ctx echo.Context) error
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
	companyIdStr := ctx.Param("company_id")
	limitStr := ctx.QueryParam("limit")
	offsetStr := ctx.QueryParam("offset")

	companyId, err := strconv.ParseUint(companyIdStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "company_id is bad request"})
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = DEFAULT_TASK_LIMIT
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = DEFAULT_TASK_OFFSET
	}

	tasks, err := c.taskUseCase.GetTasks(uint(companyId), limit, offset)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}

		slog.Info(fmt.Sprintf("error GetTasks: %v", err))
		return ctx.JSON(http.StatusInternalServerError, nil)
	}
	return ctx.JSON(http.StatusOK, response.NewGetTasksResponseBody(tasks))
}

func (c *taskController) GetTask(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("task_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}
	companyId, err := strconv.ParseUint(ctx.Param("company_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	task, err := c.taskUseCase.GetTask(uint(companyId), uint(id))
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}

		slog.Info(fmt.Sprintf("error GetTask: %v", err))
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get task"})
	}

	return ctx.JSON(http.StatusOK, response.NewGetTaskResponseBody(task))
}

func (c *taskController) CreateTask(ctx echo.Context) error {
	companyId, err := strconv.ParseUint(ctx.Param("company_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	requestBody := &request.CreateTaskRequestBody{}
	if err := ctx.Bind(requestBody); err != nil {
		return err
	}

	// Validation
	if err := validate.Struct(requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task := request.NewTaskFromCreateTaskRequestBody(uint(companyId), requestBody)
	task, err = c.taskUseCase.CreateTask(task)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create task"})
	}

	return ctx.JSON(http.StatusCreated, response.NewCreateTaskResponseBody(task))
}

func (c *taskController) UpdateTask(ctx echo.Context) error {
	taskId, err := strconv.ParseUint(ctx.Param("task_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	companyId, err := strconv.ParseUint(ctx.Param("company_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	requestBody := &request.UpdateTaskRequestBody{}
	if err := ctx.Bind(requestBody); err != nil {
		return err
	}

	// Validation
	if err := validate.Struct(requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task := request.NewTaskFromUpdateTaskRequestBody(uint(taskId), uint(companyId), requestBody)
	task, err = c.taskUseCase.UpdateTask(uint(companyId), uint(taskId), task)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}

		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, response.NewUpdateTaskResponseBody(task))
}

func (c *taskController) DeleteTask(ctx echo.Context) error {
	taskId, err := strconv.ParseUint(ctx.Param("task_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	companyId, err := strconv.ParseUint(ctx.Param("company_id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	err = c.taskUseCase.DeleteTask(uint(companyId), uint(taskId))
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}

		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusOK, nil)
}
