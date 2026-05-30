package handler

import (
	"fmt"
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

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
			"key":         v.Key,
			"title":       v.Title,
			"description": v.Description,
			"priority":    v.Priority,
			"startDate":   v.StartDate,
			"endDate":     v.EndDate,
			"allDay":      v.AllDay,
			"assignees":   v.Assigns,
		})
	}

	return ResOk(c, fiber.StatusOK, taskMap, &total, &opts)
}

func (h *TaskHandler) GetByKey(c *fiber.Ctx) error {
	key := c.Params("key")

	task, err := h.svc.GetByKey(c.Context(), key)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch task", err.Error())
	}

	return ResOk(c, fiber.StatusOK, task, nil, nil)
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

func (h *TaskHandler) CreateListStatus(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("project_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req []domain.CreateListTaskStatusReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.svc.CreateListStatus(c.Context(), projectID, req); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create task statuses", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
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

func (h *TaskHandler) DeleteStatus(c *fiber.Ctx) error {
	statusID, err := strconv.ParseInt(c.Params("status_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	if err := h.svc.DeleteStatus(c.Context(), statusID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete task status", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *TaskHandler) Create(c *fiber.Ctx) error {
	var req domain.TaskReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	userID := getCallerID(c)

	task, err := h.svc.Create(c.Context(), &req, userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create task", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, task, nil, nil)
}

func (h *TaskHandler) Update(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.UpdateTaskReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	fmt.Printf("Received update request: %+v\n", req)

	task, err := h.svc.Update(c.Context(), &req, taskID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update task", err.Error())
	}

	return ResOk(c, fiber.StatusOK, task, nil, nil)
}

func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	if err := h.svc.Delete(c.Context(), taskID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete task", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
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

func (h *TaskHandler) ListByToday(c *fiber.Ctx) error {

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	userID := getCallerID(c)
	orgID := getCallerOrgID(c)

	tasks, total, err := h.svc.ListByToday(c.Context(), opts, userID, orgID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch tasks", err.Error())
	}

	return ResOk(c, fiber.StatusOK, tasks, &total, &opts)
}

func (h *TaskHandler) ListOverdue(c *fiber.Ctx) error {
	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	userID := getCallerID(c)
	orgID := getCallerOrgID(c)

	tasks, total, err := h.svc.ListOverdue(c.Context(), opts, userID, orgID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch tasks", err.Error())
	}

	return ResOk(c, fiber.StatusOK, tasks, &total, &opts)
}

func (h *TaskHandler) ListSubTasks(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	subtasks, err := h.svc.ListSubTasks(c.Context(), taskID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch subtasks", err.Error())
	}

	return ResOk(c, fiber.StatusOK, subtasks, nil, nil)
}

func (h *TaskHandler) CreateSubTask(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.SubTaskReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}
	req.TaskID = taskID

	userID := getCallerID(c)

	subtask, err := h.svc.CreateSubTask(c.Context(), &req, userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create subtask", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, subtask, nil, nil)
}

func (h *TaskHandler) UpdateSubTask(c *fiber.Ctx) error {
	subtaskID, err := strconv.ParseInt(c.Params("subtask_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.UpdateSubTaskReq
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	subtask, err := h.svc.UpdateSubTask(c.Context(), &req, subtaskID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update subtask", err.Error())
	}

	return ResOk(c, fiber.StatusOK, subtask, nil, nil)
}

func (h *TaskHandler) DeleteSubTask(c *fiber.Ctx) error {
	subtaskID, err := strconv.ParseInt(c.Params("subtask_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	if err := h.svc.DeleteSubTask(c.Context(), subtaskID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete subtask", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *TaskHandler) ReorderSubTask(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var req domain.ReqReorderSubTask
	if err := c.BodyParser(&req); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request body", err.Error())
	}

	if err := h.svc.ReorderSubTask(c.Context(), &req, taskID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to reorder subtasks", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}


func (h *TaskHandler) GetAttachments(c *fiber.Ctx) error {
	taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid task id", err.Error())
	}
	items, err := h.svc.GetAttachments(c.Context(), taskID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to get attachments", err.Error())
	}
	return ResOk(c, fiber.StatusOK, items, nil, nil)
}

func (h *TaskHandler) DeleteAttachment(c *fiber.Ctx) error {
	attachmentID, err := strconv.ParseInt(c.Params("attachment_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid attachment id", err.Error())
	}
	if err := h.svc.DeleteAttachment(c.Context(), attachmentID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete attachment", err.Error())
	}
	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *TaskHandler) DownloadAttachment(c *fiber.Ctx) error {
	attachmentID, err := strconv.ParseInt(c.Params("attachment_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid attachment id", err.Error())
	}
	res, err := h.svc.DownloadAttachment(c.Context(), attachmentID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to download attachment", err.Error())
	}
	defer res.Reader.Close()

	c.Set("Content-Type", res.MimeType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, res.Filename))
	c.Set("Content-Length", strconv.FormatInt(res.Size, 10))
	return c.SendStream(res.Reader, int(res.Size))
}

func (h *TaskHandler) CreateAttachments(c *fiber.Ctx) error {
		taskID, err := strconv.ParseInt(c.Params("task_id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid task id", err.Error())
	}

	userID := getCallerID(c)

	form, err := c.MultipartForm()
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid multipart form", err.Error())
	}

	fileHeaders := form.File["files"]
	if len(fileHeaders) == 0 {
		return ResError(c, fiber.StatusBadRequest, "files are required", "at least one file must be provided")
	}

	const maxSize = 10 << 20 // 10MB per file

	items := make([]domain.UploadFileItem, 0, len(fileHeaders))
	closers := make([]func(), 0, len(fileHeaders))
	defer func() {
		for _, close := range closers {
			close()
		}
	}()

	for _, fh := range fileHeaders {
		if fh.Size > maxSize {
			return ResError(c, fiber.StatusBadRequest, "file too large", fh.Filename+" exceeds 10MB limit")
		}

		file, err := fh.Open()
		if err != nil {
			return ResError(c, fiber.StatusInternalServerError, "failed to open file", err.Error())
		}
		closers = append(closers, func() { file.Close() })

		items = append(items, domain.UploadFileItem{
			File:        file,
			Size:        fh.Size,
			ContentType: fh.Header.Get("Content-Type"),
			Filename:    fh.Filename,
		})
	}

	results, err := h.svc.CreateAttachments(c.Context(), items, taskID, userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to upload files", err.Error())
	}

	return ResOk(c, fiber.StatusOK, results, nil, nil)
}