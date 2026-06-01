package handler

import (
	"strconv"

	"aprilpollo/internal/adapters/hub"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type TaskCommentWSHandler struct {
	hub *hub.Hub
}

func NewTaskCommentWSHandler(h *hub.Hub) *TaskCommentWSHandler {
	return &TaskCommentWSHandler{hub: h}
}

// RequireUpgrade rejects non-WebSocket requests before the upgrade.
func (h *TaskCommentWSHandler) RequireUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Handle is the WebSocket handler. It registers the connection in the hub,
// keeps it alive by reading frames, and unregisters on disconnect.
func (h *TaskCommentWSHandler) Handle(c *websocket.Conn) {
	taskID, err := strconv.ParseInt(c.Params("taskID"), 10, 64)
	if err != nil {
		_ = c.Close()
		return
	}

	h.hub.Register(taskID, c)
	defer h.hub.Unregister(taskID, c)

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
