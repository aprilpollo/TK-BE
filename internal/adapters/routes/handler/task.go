package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"
	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	svc input.TaskService
}

func NewTaskHandler(svc input.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) List(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}
	statusID, err := strconv.ParseInt(c.Params("status_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	tasks, total, err := h.svc.List(c.Context(), opts, projectID, statusID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch tasks", err.Error())
	}
	taskMap := make([]map[string]interface{}, 0)
	for _, v := range tasks {
		taskMap = append(taskMap, map[string]interface{}{
			"id":          v.ID,
			"columnId":    v.Status.UUID,
			"title":       v.Title,
			"description": v.Description,
			"priority":    v.Priority,
			"dueDate":     v.DueDate,
			"subtasks":    v.ParentID,
		})
	}

	return ResOk(c, fiber.StatusOK, taskMap, &total, &opts)
}

func (h *TaskHandler) ListPriority(c *fiber.Ctx) error {
	priorities, err := h.svc.ListPriority(c.Context())
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch task priorities", err.Error())
	}

	return ResOk(c, fiber.StatusOK, priorities, nil, nil)
}

func (h *TaskHandler) ListStatus(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	statuses, err := h.svc.ListStatus(c.Context(), opts, projectID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch task statuses", err.Error())
	}

	return ResOk(c, fiber.StatusOK, statuses, nil, nil)
}

func (h *TaskHandler) CreateStatus(c *fiber.Ctx) error {
	var req domain.CreateTaskStatusReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	status, err := h.svc.CreateStatus(c.Context(), &req)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create task status", err.Error())
	}

	return ResOk(c, fiber.StatusOK, status, nil, nil)
}

func (h *TaskHandler) UpdateStatus(c *fiber.Ctx) error {
	statusID, err := strconv.ParseInt(c.Params("status_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.UpdateTaskStatusReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	status, err := h.svc.UpdateStatus(c.Context(), &req, statusID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update task status", err.Error())
	}

	return ResOk(c, fiber.StatusOK, status, nil, nil)
}

func (h *TaskHandler) Create(c *fiber.Ctx) error {
	var req domain.TaskReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	task, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create task", err.Error())
	}

	return ResOk(c, fiber.StatusOK, task, nil, nil)
}

func (h *TaskHandler) ReorderStatus(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.ReqReorderTaskStatus
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.svc.ReorderStatus(c.Context(), &req, projectID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to reorder task statuses", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *TaskHandler) ReorderTask(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.ReqReorderTask
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.svc.ReorderTask(c.Context(), &req, projectID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to reorder tasks", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}
