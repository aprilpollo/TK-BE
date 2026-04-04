package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct {
	svc input.OrganizationService
}

func NewOrganizationHandler(svc input.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

// GET /api/v1/organizations
func (h *OrganizationHandler) Gets(c *fiber.Ctx) error {
	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	orgs, total, err := h.svc.List(c.Context(), opts)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch organizations", err.Error())
	}

	return ResOk(c, fiber.StatusOK, orgs, &total, &opts)
}

// GET /api/v1/organizations/:id
func (h *OrganizationHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	org, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch organization", err.Error())
	}
	if org == nil {
		return ResError(c, fiber.StatusNotFound, "organization not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, org, nil, nil)
}

// POST /api/v1/organizations
func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var body domain.CreateOrganizationReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	org, err := h.svc.Create(c.Context(), &body, getCallerID(c))
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create organization", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, org, nil, nil)
}

// PUT /api/v1/organizations/:id
func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var body domain.UpdateOrganizationReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	org, err := h.svc.Update(c.Context(), id, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update organization", err.Error())
	}

	return ResOk(c, fiber.StatusOK, org, nil, nil)
}

// DELETE /api/v1/organizations/:id
func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete organization", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GET /api/v1/organizations/:id/members
func (h *OrganizationHandler) GetMembers(c *fiber.Ctx) error {
	orgID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	members, total, err := h.svc.ListMembers(c.Context(), orgID, opts)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch members", err.Error())
	}

	return ResOk(c, fiber.StatusOK, members, &total, &opts)
}

// POST /api/v1/organizations/:id/members
func (h *OrganizationHandler) InviteMember(c *fiber.Ctx) error {
	orgID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	var body domain.InviteMemberReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	member, err := h.svc.InviteMember(c.Context(), orgID, &body, getCallerID(c))
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to invite member", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, member, nil, nil)
}

// PUT /api/v1/organizations/:id/members/:memberID
func (h *OrganizationHandler) UpdateMember(c *fiber.Ctx) error {
	orgID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	memberID, err := strconv.ParseInt(c.Params("memberID"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid member id", err.Error())
	}

	var body domain.UpdateMemberReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	member, err := h.svc.UpdateMember(c.Context(), orgID, memberID, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update member", err.Error())
	}

	return ResOk(c, fiber.StatusOK, member, nil, nil)
}

// DELETE /api/v1/organizations/:id/members/:memberID
func (h *OrganizationHandler) RemoveMember(c *fiber.Ctx) error {
	orgID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	memberID, err := strconv.ParseInt(c.Params("memberID"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid member id", err.Error())
	}

	if err := h.svc.RemoveMember(c.Context(), orgID, memberID); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to remove member", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
