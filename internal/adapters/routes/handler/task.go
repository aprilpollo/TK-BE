package handler

import (
	// "strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	// "aprilpollo/internal/pkg/query"
	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	svc input.TaskService
}

func NewTaskHandler(svc input.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) ListPriority(c *fiber.Ctx) error {
	priorities, err := h.svc.ListPriority(c.Context())
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch task priorities", err.Error())
	}

	return ResOk(c, fiber.StatusOK, priorities, nil, nil)
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