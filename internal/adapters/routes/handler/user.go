package handler

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc    input.UserService
	svcOrg input.OrganizationService
}

func NewUserHandler(svc input.UserService, svcOrg input.OrganizationService) *UserHandler {
	return &UserHandler{svc: svc, svcOrg: svcOrg}
}

// GET /api/v1/users
func (h *UserHandler) Gets(c *fiber.Ctx) error {
	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	users, total, err := h.svc.List(c.Context(), opts)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch users", err.Error())
	}

	return ResOk(c, fiber.StatusOK, users, &total, &opts)
}

// GET /api/v1/users/me
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userID := getCallerID(c)

	user, err := h.svc.GetByID(c.Context(), userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch user", err.Error())
	}
	if user == nil {
		return ResError(c, fiber.StatusNotFound, "user not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, user, nil, nil)
}

// GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid id", err.Error())
	}

	user, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch user", err.Error())
	}
	if user == nil {
		return ResError(c, fiber.StatusNotFound, "user not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, user, nil, nil)
}

// GET /api/v1/users/me/organizations
func (h *UserHandler) GetMyOrganizations(c *fiber.Ctx) error {
	userID := getCallerID(c)

	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	orgs, total, err := h.svc.ListMyOrganizations(c.Context(), userID, opts)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch organizations", err.Error())
	}

	return ResOk(c, fiber.StatusOK, orgs, &total, &opts)
}

// GET /api/v1/users/me/organizations/permissions
func (h *UserHandler) GetMyPrimaryOrgPermissions(c *fiber.Ctx) error {
	userID := getCallerID(c)

	perms, err := h.svc.GetMyPrimaryOrgPermissions(c.Context(), userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch permissions", err.Error())
	}
	if perms == nil {
		return ResError(c, fiber.StatusNotFound, "no primary organization found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, perms, nil, nil)
}

// PUT /api/v1/users/me
func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	userID := getCallerID(c)

	var body domain.UpdateUserReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	user, err := h.svc.Update(c.Context(), userID, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update user", err.Error())
	}

	return ResOk(c, fiber.StatusOK, user, nil, nil)
}

func (h *UserHandler) UpdatePrimaryOrganization(c *fiber.Ctx) error {
	userID := getCallerID(c)
	orgID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid organization id", err.Error())
	}

	err = h.svcOrg.UpdatePrimary(c.Context(), orgID, userID)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update primary organization", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}

func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error {
	userID := getCallerID(c)

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "avatar file is required", err.Error())
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

	req := &domain.AvatarUploadReq{
		File:        file,
		Size:        fileHeader.Size,
		ContentType: contentType,
		Filename:    fileHeader.Filename,
	}

	if err := h.svc.UpdateAvatar(c.Context(), userID, req); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update avatar", err.Error())
	}

	return ResOk(c, fiber.StatusOK, nil, nil, nil)
}
