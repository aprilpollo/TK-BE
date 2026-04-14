package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

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

func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var body domain.CreateProjectReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	project, err := h.svc.Create(c.Context(), &body)
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

	err = h.svc.Update(c.Context(), id, &body, orgId)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update project", err.Error())
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
