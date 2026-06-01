package routes

import (
	"aprilpollo/internal/adapters/routes/handler"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterOauthRoutes(app *fiber.App, h *handler.OauthHandler) {
	api := app.Group("/api/v1")

	oauth := api.Group("/auth")
	oauth.Post("/basic-login", h.BasicLogin)
	oauth.Post("/social-login", h.SocialLogin)
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
	projects.Get("/:id/task-summary", h.GetTaskSummary)
	projects.Get("/:id/chart", h.GetTaskVelocityChart)
	projects.Get("/:id/members", h.GetProjectMembers)
	projects.Get("/:id/task-deadlines", h.GetUpcomingDeadlines)

	projects.Post("/", h.Create)
	projects.Post("/:id/logo", h.UpdateLogo)

	projects.Put("/:id", h.Update)
	projects.Put("/:id/notification", h.UpdateNotificationSettings)

	projects.Delete("/:id", h.Delete)
}

func RegisterTaskRoutes(app *fiber.App, h *handler.TaskHandler, ch *handler.TaskCommentHandler, wsHandler *handler.TaskCommentWSHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	tasks := api.Group("/tasks", jwtMiddleware, orgMiddleware)
	tasks.Get("/priorities", h.ListPriority)
	tasks.Get("/statuses/:projectID", h.ListStatus)
	tasks.Get("/me/today", h.ListByToday)
	tasks.Get("/me/overdue", h.ListOverdue)
	tasks.Get("/key/:key", h.GetByKey)

	tasks.Get("/:taskID/attachments", h.GetAttachments)
	tasks.Post("/:taskID/attachments", h.CreateAttachments)
	tasks.Delete("/:taskID/attachments/:attachmentID", h.DeleteAttachment)

	// comment routes registered before the /:projectID/:statusID wildcard
	// so Fiber's radix tree prefers the static "comments" segment
	tasks.Get("/:taskID/comments", ch.List)
	tasks.Post("/:taskID/comments/upload", ch.UploadFile)
	tasks.Post("/:taskID/comments", ch.Create)
	tasks.Put("/:taskID/comments/:commentID", ch.Update)
	tasks.Delete("/:taskID/comments/:commentID", ch.Delete)

	// real-time comments via WebSocket: ws://.../api/v1/tasks/:taskID/comments/live
	tasks.Get("/:taskID/comments/live", wsHandler.RequireUpgrade, websocket.New(wsHandler.Handle))

	// subtask routes
	tasks.Get("/:taskID/subtasks", h.ListSubTasks)
	tasks.Post("/:taskID/subtasks", h.CreateSubTask)
	tasks.Put("/:taskID/subtasks/reorder", h.ReorderSubTask)
	tasks.Put("/:taskID/subtasks/:subtaskID", h.UpdateSubTask)
	tasks.Delete("/:taskID/subtasks/:subtaskID", h.DeleteSubTask)

	tasks.Get("/:projectID/:statusID", h.List)

	tasks.Post("/", h.Create)
	tasks.Post("/statuses", h.CreateStatus)
	tasks.Post("/statuses/list/:projectID", h.CreateListStatus)

	tasks.Put("/:taskID", h.Update)
	tasks.Put("/statuses/reorder/:projectID", h.ReorderStatus)
	tasks.Put("/statuses/:statusID", h.UpdateStatus)
	tasks.Put("/reorder/:projectID", h.ReorderTask)

	tasks.Delete("/:taskID", h.Delete)
	tasks.Delete("/statuses/:statusID", h.DeleteStatus)
}

func RegisterCalendarRoutes(app *fiber.App, h *handler.CalendarHandler, jwtMiddleware fiber.Handler, orgMiddleware fiber.Handler) {
	api := app.Group("/api/v1")

	calendar := api.Group("/calendar", jwtMiddleware, orgMiddleware)
	calendar.Get("/priorities", h.ListPriority)
	calendar.Get("/statuses/:projectID", h.ListStatus)
	calendar.Get("/:projectID", h.List)

}
