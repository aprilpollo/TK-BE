package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

	"github.com/gofiber/fiber/v2"
)

type TaskCommentHandler struct {
	svc input.TaskCommentService
}

func NewTaskCommentHandler(svc input.TaskCommentService) *TaskCommentHandler {
	return &TaskCommentHandler{svc: svc}
}

func (h *TaskCommentHandler) List(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid task id", err.Error())
	}

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	comments, total, err := h.svc.List(c.Context(), opts, taskID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch comments", err.Error())
	}

	return ResOk(c, fiber.StatusOK, comments, &total, &opts)
}

func (h *TaskCommentHandler) Create(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid task id", err.Error())
	}

	var req domain.CreateTaskCommentReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	userID := getCallerID(c)

	comment, err := h.svc.Create(c.Context(), &req, taskID, userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create comment", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, comment, nil, nil)
}

func (h *TaskCommentHandler) Update(c *fiber.Ctx) error {
	commentID, err := strconv.ParseInt(c.Params("comment_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid comment id", err.Error())
	}

	var req domain.UpdateTaskCommentReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	comment, err := h.svc.Update(c.Context(), &req, commentID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update comment", err.Error())
	}

	return ResOk(c, fiber.StatusOK, comment, nil, nil)
}

func (h *TaskCommentHandler) Delete(c *fiber.Ctx) error {
	commentID, err := strconv.ParseInt(c.Params("comment_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid comment id", err.Error())
	}

	if err := h.svc.Delete(c.Context(), commentID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete comment", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}
