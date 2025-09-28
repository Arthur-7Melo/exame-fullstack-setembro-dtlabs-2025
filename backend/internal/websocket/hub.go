package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	logger "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	redis      *redis.Client
	mu         sync.RWMutex
	shutdown   chan struct{}
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redisClient,
		shutdown:   make(chan struct{}),
	}
}

func safeClose(ch chan []byte) {
	defer func() {
		_ = recover()
	}()
	if ch != nil {
		close(ch)
	}
}

func (h *Hub) Run() {
	logger.Logger.Info("WebSocket Hub started")
	go h.listenRedis()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if existing, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				safeClose(existing.Send)
				existing.Conn.Close()
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			logger.Logger.Info("Client registered", "user_id", client.UserID, "total_clients", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if existing, ok := h.clients[client.UserID]; ok && existing == client {
				delete(h.clients, client.UserID)
				safeClose(client.Send)
				client.Conn.Close()
				logger.Logger.Info("Client unregistered", "user_id", client.UserID, "total_clients", len(h.clients))
			}
			h.mu.Unlock()

		case <-h.shutdown:
			h.mu.Lock()
			for userID, client := range h.clients {
				safeClose(client.Send)
				client.Conn.Close()
				delete(h.clients, userID)
			}
			h.mu.Unlock()
			logger.Logger.Info("WebSocket Hub stopped")
			return
		}
	}
}

func (h *Hub) listenRedis() {
	logger.Logger.Info("Starting Redis listener for notifications")

	for {
		ctx, cancel := context.WithCancel(context.Background())
		pubsub := h.redis.PSubscribe(ctx, "notifications:*")

		_, err := pubsub.Receive(ctx)
		if err != nil {
			logger.Logger.Error("Failed to subscribe to Redis", "error", err)
			pubsub.Close()
			cancel()
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Logger.Info("Successfully subscribed to Redis notifications channel")
		channel := pubsub.Channel()

	redisLoop:
		for {
			select {
			case msg, ok := <-channel:
				if !ok {
					logger.Logger.Warn("Redis channel closed, reconnecting...")
					break redisLoop
				}

				logger.Logger.Debug("Received message from Redis", "channel", msg.Channel, "payload", msg.Payload)

				var notification map[string]interface{}
				if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
					logger.Logger.Error("Error decoding notification", "error", err)
					continue
				}

				var userID string
				switch v := notification["user_id"].(type) {
				case string:
					userID = v
				case float64:
					userID = strconv.FormatFloat(v, 'f', -1, 64)
				default:
					logger.Logger.Warn("Notification missing/invalid user_id", "payload", msg.Payload)
					continue
				}

				h.mu.RLock()
				client, exists := h.clients[userID]
				h.mu.RUnlock()

				if exists && client != nil {
					select {
					case client.Send <- []byte(msg.Payload):
						logger.Logger.Debug("Notification sent to client", "user_id", userID)
					default:
						logger.Logger.Warn("Client channel full, disconnecting", "user_id", userID)
						h.unregister <- client
					}
				} else {
					logger.Logger.Debug("No client found for user", "user_id", userID)
				}

			case <-h.shutdown:
				pubsub.Close()
				cancel()
				return
			}
		}

		pubsub.Close()
		cancel()
		time.Sleep(1 * time.Second)
	}
}

func (h *Hub) HandleWebSocket(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id é obrigatório"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Error("Error upgrading to WebSocket", "error", err, "user_id", userID)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *Hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Logger.Error("WebSocket read error", "error", err, "user_id", client.UserID)
			}
			break
		}
	}
}

func (h *Hub) writePump(client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Logger.Error("NextWriter error", "error", err, "user_id", client.UserID)
				h.unregister <- client
				return
			}
			if _, err := w.Write(message); err != nil {
				logger.Logger.Error("Write message error", "error", err, "user_id", client.UserID)
			}

			n := len(client.Send)
			for i := 0; i < n; i++ {
				if msg := <-client.Send; msg != nil {
					if _, err := w.Write(msg); err != nil {
						logger.Logger.Error("Write queued message error", "error", err, "user_id", client.UserID)
					}
				}
			}

			if err := w.Close(); err != nil {
				logger.Logger.Error("Writer close error", "error", err, "user_id", client.UserID)
				h.unregister <- client
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Logger.Error("Ping error", "error", err, "user_id", client.UserID)
				h.unregister <- client
				return
			}
		}
	}
}

func (h *Hub) Stop() {
	select {
	case <-h.shutdown:
	default:
		close(h.shutdown)
	}
}
