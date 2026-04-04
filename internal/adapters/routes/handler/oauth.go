package handler

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"

	"github.com/gofiber/fiber/v2"
)

type OauthHandler struct {
	svc input.OauthService
}

func NewOauthHandler(svc input.OauthService) *OauthHandler {
	return &OauthHandler{svc: svc}
}

func (h *OauthHandler) BasicLogin(c *fiber.Ctx) error {
	var body domain.BasicLogin
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	token, err := h.svc.BasicLogin(c.Context(), &body)
	if err != nil {
		return ResError(c, fiber.StatusUnauthorized, "unauthorized", err.Error())
	}

	return ResOk(c, fiber.StatusOK, fiber.Map{"token": token}, nil, nil)
}

func (h *OauthHandler) SocialLogin(c *fiber.Ctx) error {
	var body domain.SocialLogin
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	token, err := h.svc.SocialLogin(c.Context(), &body)
	if err != nil {
		return ResError(c, fiber.StatusUnauthorized, "unauthorized", err.Error())
	}

	return ResOk(c, fiber.StatusOK, fiber.Map{"token": token}, nil, nil)
}

func (h *OauthHandler) Register(c *fiber.Ctx) error {
	return nil
}

func (h *OauthHandler) RefreshToken(c *fiber.Ctx) error {
	return nil
}
