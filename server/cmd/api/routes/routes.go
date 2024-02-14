package routes

import (
	lg "gotrader/logger"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func init() {
	logger = lg.CreateCustomLogger("server/routes")
}

var logger *log.Logger

type request struct {
	Exchange string `json:"exchange,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (m *Manager) HandleRoutes(w http.ResponseWriter, r *http.Request) {
	m.writeJsonResponse(w, http.StatusOK, jsonResponse{
		Success: true,
		Message: "API is working",
		Data:    nil,
	})
}

func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		conn.WriteJSON(jsonResponse{
			Success: false,
			Message: "Failed to upgrade",
			Data:    nil,
		})
		return
	}

	client := NewWebsocketClient(conn, m)
	m.AddClient(client)
	conn.WriteJSON(jsonResponse{
		Success: true,
		Message: "Connected to server",
	})

	go client.ReadMessages()
}
