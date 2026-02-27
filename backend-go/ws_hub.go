package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsClient struct {
	userID int64
	send   chan []byte
}

type wsBroadcast struct {
	userID int64
	msg    []byte
}

type wsHub struct {
	clients    map[*wsClient]bool
	byUser     map[int64]map[*wsClient]bool
	register   chan *wsClient
	unregister chan *wsClient
	broadcast  chan wsBroadcast
	mu         sync.RWMutex
}

var globalWSHub *wsHub

func initWSHub() {
	globalWSHub = &wsHub{
		clients:    make(map[*wsClient]bool),
		byUser:     make(map[int64]map[*wsClient]bool),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		broadcast:  make(chan wsBroadcast, 64),
	}
	go globalWSHub.run()
}

func (h *wsHub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			if h.byUser[c.userID] == nil {
				h.byUser[c.userID] = make(map[*wsClient]bool)
			}
			h.byUser[c.userID][c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				if m := h.byUser[c.userID]; m != nil {
					delete(m, c)
					if len(m) == 0 {
						delete(h.byUser, c.userID)
					}
				}
				close(c.send)
			}
			h.mu.Unlock()
		case b := <-h.broadcast:
			h.mu.RLock()
			for cl := range h.byUser[b.userID] {
				select {
				case cl.send <- b.msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToUser sends a JSON message to all connections for the given user.
func BroadcastToUser(userID int64, event string, payload interface{}) {
	if globalWSHub == nil {
		return
	}
	out := map[string]interface{}{"event": event}
	if payload != nil {
		out["data"] = payload
	}
	raw, err := json.Marshal(out)
	if err != nil {
		return
	}
	select {
	case globalWSHub.broadcast <- wsBroadcast{userID: userID, msg: raw}:
	default:
	}
}

func handleWS(c *gin.Context) {
	uid := getUserID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &wsClient{userID: uid, send: make(chan []byte, 64)}
	globalWSHub.register <- client
	defer func() { globalWSHub.unregister <- client }()
	go func() {
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}()
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
