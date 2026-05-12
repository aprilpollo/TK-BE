package routes

import (
	"aprilpollo/internal/adapters/routes/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterOauthRoutes(app *fiber.App, h *handler.OauthHandler) {
	api := app.Group("/api/v1")

	oauth := api.Group("/auth")
	oauth.Post("/basiclogin", h.BasicLogin)
	oauth.Post("/sociallogin", h.SocialLogin)
}

func RegisterUserRoutes(app *fiber.App, h *handler.UserHandler, jwtMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	users := api.Group("/users", jwtMiddleware)
	users.Get("/me", h.GetMe)
	users.Get("/me/organizations", h.GetMyOrganizations)
	users.Get("/me/organizations/permissions", h.GetMyPrimaryOrgPermissions)

	users.Post("/me/avatar", h.UpdateAvatar)

	users.Put("/me", h.UpdateMe)
	users.Put("/me/organizations/primary/:id", h.UpdatePrimaryOrganization)
}

func RegisterOrganizationRoutes(app *fiber.App, h *handler.OrganizationHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	orgs := api.Group("/organizations", jwtMiddleware)
	orgs.Get("/", h.Gets)
	orgs.Get("/members", orgMiddleware, h.GetMembers)
	orgs.Get("/:id", h.GetByID)

	orgs.Post("/", h.Create)
	orgs.Post("/:id/members", h.InviteMember)

	orgs.Put("/:id", h.Update)
	orgs.Put("/:id/members/:memberID", h.UpdateMember)

	orgs.Delete("/:id", h.Delete)
	orgs.Delete("/:id/members/:memberID", h.RemoveMember)
}

func RegisterProjectRoutes(app *fiber.App, h *handler.ProjectHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	projects := api.Group("/projects", jwtMiddleware, orgMiddleware)
	projects.Get("/", h.Gets)
	projects.Get("/statuses", h.GetStatuses)
	projects.Get("/key/:key", h.GetByKey)
	projects.Get("/:id", h.GetByID)
	projects.Get("/:id/notification", h.GetNotificationSettings)
	projects.Get("/:id/task_summary", h.GetTaskSummary)
	projects.Get("/:id/chart", h.GetTaskVelocityChart)
	projects.Get("/:id/members", h.GetProjectMembers)
	projects.Get("/:id/task_deadlines", h.GetUpcomingDeadlines)

	projects.Post("/", h.Create)
	projects.Post("/:id/logo", h.UpdateLogo)

	projects.Put("/:id", h.Update)
	projects.Put("/:id/notification", h.UpdateNotificationSettings)

	projects.Delete("/:id", h.Delete)
}

func RegisterTaskRoutes(app *fiber.App, h *handler.TaskHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	tasks := api.Group("/tasks", jwtMiddleware, orgMiddleware)
	tasks.Get("/priorities", h.ListPriority)
	tasks.Get("/statuses/:project_id", h.ListStatus)
	tasks.Get("/me/today", h.ListByToday)
	tasks.Get("/me/overdue", h.ListOverdue)
	tasks.Get("/:project_id/:status_id", h.List)

	tasks.Post("/", h.Create)
	tasks.Post("/statuses", h.CreateStatus)
	tasks.Post("/statuses/list/:project_id", h.CreateListStatus)

	tasks.Put("/:task_id", h.Update)
	tasks.Put("/statuses/reorder/:project_id", h.ReorderStatus)
	tasks.Put("/statuses/:status_id", h.UpdateStatus)

	tasks.Put("/reorder/:project_id", h.ReorderTask)

	tasks.Delete("/:task_id", h.Delete)
	tasks.Delete("/statuses/:status_id", h.DeleteStatus)
}

func RegisterCalendarRoutes(app *fiber.App, h *handler.CalendarHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	calendar := api.Group("/calendar", jwtMiddleware, orgMiddleware)
	calendar.Get("/priorities", h.ListPriority)
	calendar.Get("/statuses/:project_id", h.ListStatus)
	calendar.Get("/:project_id", h.List)

}
