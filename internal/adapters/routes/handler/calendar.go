package handler

import (
	"strconv"

	// "aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type CalendarHandler struct {
	svc input.CalendarService
}

func NewCalendarHandler(svc input.CalendarService) *CalendarHandler {
	return &CalendarHandler{svc: svc}
}

func (h *CalendarHandler) List(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	tasks, total, err := h.svc.List(c.Context(), opts, projectID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch tasks", err.Error())
	}
	taskMap := make([]map[string]interface{}, 0)
	for _, v := range tasks {
		taskMap = append(taskMap, map[string]interface{}{
			"id":          v.Key,
			"title":       v.Title,
			"description": v.Description,
			"start":       v.StartDate,
			"end":         v.EndDate,
			"allDay":      v.AllDay,
			"category":    "task",
			"priority":    v.Priority,
			"assignees":   v.Assigns,
			"status":      v.Status,
		})
	}

	return ResOk(c, fiber.StatusOK, taskMap, &total, &opts)
}

func (h *CalendarHandler) ListPriority(c *fiber.Ctx) error {
	priorities, err := h.svc.ListPriority(c.Context())
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch task priorities", err.Error())
	}

	return ResOk(c, fiber.StatusOK, priorities, nil, nil)
}

func (h *CalendarHandler) ListStatus(c *fiber.Ctx) error {
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
