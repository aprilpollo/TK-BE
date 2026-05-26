package hub

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/redis/go-redis/v9"
)

// Hub manages WebSocket connections grouped by task ID.
// A single Redis PSubscribe goroutine fans messages out to the right room.
type Hub struct {
	mu    sync.RWMutex
	rooms map[int64]map[*websocket.Conn]struct{}
}

func New() *Hub {
	return &Hub{rooms: make(map[int64]map[*websocket.Conn]struct{})}
}

func (h *Hub) Register(taskID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[taskID] == nil {
		h.rooms[taskID] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[taskID][conn] = struct{}{}
}

func (h *Hub) Unregister(taskID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms[taskID], conn)
	if len(h.rooms[taskID]) == 0 {
		delete(h.rooms, taskID)
	}
}

func (h *Hub) Broadcast(taskID int64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.rooms[taskID] {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("hub: ws write error task=%d err=%v", taskID, err)
		}
	}
}

// StartSubscriber starts a background goroutine that subscribes to
// "task:comments:*" on Redis and broadcasts payloads to the matching room.
func (h *Hub) StartSubscriber(ctx context.Context, redisClient *redis.Client) {
	pubsub := redisClient.PSubscribe(ctx, "task:comments:*")
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				// channel format: "task:comments:{taskID}"
				parts := strings.Split(msg.Channel, ":")
				if len(parts) != 3 {
					continue
				}
				taskID, err := strconv.ParseInt(parts[2], 10, 64)
				if err != nil {
					continue
				}
				h.Broadcast(taskID, []byte(msg.Payload))
			}
		}
	}()
}
