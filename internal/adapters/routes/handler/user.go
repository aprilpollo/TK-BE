package handler

import (
	"strconv"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc input.UserService
}

func NewUserHandler(svc input.UserService) *UserHandler {
	return &UserHandler{svc: svc}
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
