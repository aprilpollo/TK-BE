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

func (h *TaskCommentHandler) UploadFile(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid task id", err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid multipart form", err.Error())
	}

	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		return ResError(c, fiber.StatusBadRequest, "files are required", "at least one file must be provided")
	}

	const maxSize = 10 << 20 // 10MB per file

	results := make([]*domain.TaskCommentFileUploadRes, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		if fh.Size > maxSize {
			return ResError(c, fiber.StatusBadRequest, "file too large", fh.Filename+" exceeds 10MB limit")
		}

		file, err := fh.Open()
		if err != nil {
			return ResError(c, fiber.StatusInternalServerError, "failed to open file", err.Error())
		}

		contentType := fh.Header.Get("Content-Type")
		result, err := h.svc.UploadFile(c.Context(), file, fh.Size, contentType, fh.Filename, taskID)
		file.Close()
		if err != nil {
			return ResError(c, fiber.StatusInternalServerError, "failed to upload file", err.Error())
		}

		results = append(results, result)
	}

	return ResOk(c, fiber.StatusOK, results, nil, nil)
}
