package handler

import (
	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/pkg/query"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	svc input.BookService
}

func NewBookHandler(svc input.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// GET /api/v1/books
func (h *BookHandler) Gets(c *fiber.Ctx) error {
	opts, err := query.Parse(c.Queries())
	if err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid query", err.Error())
	}

	books, total, err := h.svc.List(c.Context(), opts)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch books", err.Error())
	}

	return ResOk(c, fiber.StatusOK, books, &total, &opts)
}

// GET /api/v1/books/:id
func (h *BookHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	book, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to fetch book", err.Error())
	}
	if book == nil {
		return ResError(c, fiber.StatusNotFound, "book not found", "record not found")
	}

	return ResOk(c, fiber.StatusOK, book, nil, nil)
}

// POST /api/v1/books
func (h *BookHandler) Create(c *fiber.Ctx) error {
	var book domain.Book
	if err := c.BodyParser(&book); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	if err := h.svc.Create(c.Context(), &book); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to create book", err.Error())
	}

	return ResOk(c, fiber.StatusCreated, book, nil, nil)
}

// PUT /api/v1/books/:id
func (h *BookHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var body domain.UpdateBookReq
	if err := c.BodyParser(&body); err != nil {
		return ResError(c, fiber.StatusBadRequest, "invalid request", err.Error())
	}

	book, err := h.svc.Update(c.Context(), id, &body)
	if err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to update book", err.Error())
	}

	return ResOk(c, fiber.StatusOK, book, nil, nil)
}

// DELETE /api/v1/books/:id
func (h *BookHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.svc.Delete(c.Context(), id); err != nil {
		return ResError(c, fiber.StatusInternalServerError, "failed to delete book", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
