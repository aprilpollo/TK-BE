package routes

import (
	"aprilpollo/internal/adapters/routes/handler"

	"github.com/gofiber/fiber/v2"
)



func RegisterUserRoutes(app *fiber.App, h *handler.UserHandler, jwtMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	users := api.Group("/users", jwtMiddleware)
	users.Get("/", h.Gets)
	users.Get("/me", h.GetMe)
	users.Get("/me/organizations", h.GetMyOrganizations)
	users.Get("/:id", h.GetByID)
	users.Put("/me", h.UpdateMe)
}

func RegisterOauthRoutes(app *fiber.App, h *handler.OauthHandler) {
	api := app.Group("/api/v1")

	oauth := api.Group("/auth")
	oauth.Post("/basiclogin", h.BasicLogin)
	oauth.Post("/sociallogin", h.SocialLogin)
}

func RegisterOrganizationRoutes(app *fiber.App, h *handler.OrganizationHandler, jwtMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	orgs := api.Group("/organizations", jwtMiddleware)
	orgs.Get("/", h.Gets)
	orgs.Get("/:id", h.GetByID)
	orgs.Post("/", h.Create)
	orgs.Put("/:id", h.Update)
	orgs.Delete("/:id", h.Delete)

	orgs.Get("/:id/members", h.GetMembers)
	orgs.Post("/:id/members", h.InviteMember)
	orgs.Put("/:id/members/:memberID", h.UpdateMember)
	orgs.Delete("/:id/members/:memberID", h.RemoveMember)
}

func RegisterBookRoutes(app *fiber.App, h *handler.BookHandler, jwtMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	books := api.Group("/books", jwtMiddleware)
	books.Get("/", h.Gets)
	books.Get("/:id", h.GetByID)
	books.Post("/", h.Create)
	books.Put("/:id", h.Update)
	books.Delete("/:id", h.Delete)
}
