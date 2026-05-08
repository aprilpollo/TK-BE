package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	svc input.ProjectService
}

func NewProjectHandler(svc input.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) Gets(c *fiber.Ctx) error {
	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	orgId := getCallerOrgID(c)

	projects, total, err := h.svc.List(c.Context(), opts, orgId)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch projects", err.Error())
	}

	return ResOk(c, fiber.StatusOK, projects, &total, &opts)
}

func (h *ProjectHandler) GetStatuses(c *fiber.Ctx) error {
	statuses, err := h.svc.ListStatuses(c.Context())
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch project statuses", err.Error())
	}

	return ResOk(c, fiber.StatusOK, statuses, nil, nil)
}

func (h *ProjectHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}
	orgId := getCallerOrgID(c)

	project, err := h.svc.GetByID(c.Context(), id, orgId)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch project", err.Error())
	}

	if project == nil {
		return ResError(c, fiber.StatusNotFound, "project not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, project, nil, nil)

}

func (h *ProjectHandler) GetByKey(c *fiber.Ctx) error {
	key, err := uuid.Parse(c.Params("key"))
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid key", err.Error())
	}
	orgId := getCallerOrgID(c)

	project, err := h.svc.GetByKey(c.Context(), key, orgId)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch project", err.Error())
	}

	if project == nil {
		return ResError(c, fiber.StatusNotFound, "project not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, project, nil, nil)

}

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var body domain.CreateProjectReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}
	orgId := getCallerOrgID(c)

	project, err := h.svc.Create(c.Context(), orgId, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create project", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, project, nil, nil)
}

func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	var body domain.UpdateProjectReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	orgId := getCallerOrgID(c)

	err = h.svc.Update(c.Context(), id, orgId, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update project", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *ProjectHandler) UpdateLogo(c *fiber.Ctx) error {
	orgId := getCallerOrgID(c)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	fileHeader, err := c.FormFile("logo")
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "logo file is required", err.Error())
	}

	const maxSize = 2 << 20 // 2MB
	if fileHeader.Size > maxSize {
		return ResError(c, fiber.StatusBadRequest, "file too large", "maximum allowed size is 2MB")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return ResError(c, fiber.StatusBadRequest, "invalid file type", "only jpg and png are allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to open file", err.Error())
	}

	defer file.Close()

	req := &domain.LogoUploadReq{
		File:        file,
		Size:        fileHeader.Size,
		ContentType: contentType,
		Filename:    fileHeader.Filename,
	}

	if err := h.svc.UpdateLogo(c.Context(), id, orgId, req); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update logo", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	orgId := getCallerOrgID(c)

	err = h.svc.Delete(c.Context(), id, orgId)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete project", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *ProjectHandler) GetNotificationSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	settings, err := h.svc.GetNotificationSettings(c.Context(), id)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch notification settings", err.Error())
	}

	return ResOk(c, fiber.StatusOK, settings, nil, nil)
}

func (h *ProjectHandler) UpdateNotificationSettings(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var body domain.ProjectNotificationSettings
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	err = h.svc.UpdateNotificationSettings(c.Context(), id, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update notification settings", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}