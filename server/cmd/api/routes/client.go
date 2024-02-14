package routes

import (
	pb "gotrader/proto"

	"github.com/gorilla/websocket"
)

type ClientHandeller struct {
	conn    *websocket.Conn
	manager *Manager
	nseChan chan *pb.Stock
}

type clientRequest struct {
	Exchange string `json:"exchange,omitempty"`
}

type ClientList map[*ClientHandeller]bool

func NewWebsocketClient(conn *websocket.Conn, manager *Manager) *ClientHandeller {
	return &ClientHandeller{
		conn:    conn,
		manager: manager,
	}
}

func (client *ClientHandeller) SubscribeToNseChan() {
	client.nseChan = make(chan *pb.Stock)
}

func (client *ClientHandeller) SendToNseChan(stock *pb.Stock) {
	client.nseChan <- stock
}

func (client *ClientHandeller) ReadMessages() {
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			client.manager.RemoveClient(client)
			break
		}

		var req clientRequest
		err = client.conn.ReadJSON(&req)
		if err != nil {
			client.conn.WriteJSON(jsonResponse{
				Success: false,
				Message: "Failed to read message",
				Data:    nil,
			})
			client.manager.RemoveClient(client)
			break
		}
		if req.Exchange == "NSE" {
			logger.Printf("Request recieved from %v, for exchange %v", client.conn.RemoteAddr(), req.Exchange)
		}
	}
}
