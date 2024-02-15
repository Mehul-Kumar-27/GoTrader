package routes

import (
	"encoding/json"
	pb "gotrader/proto"

	//	"sync"

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

type NseClientList map[*ClientHandeller]bool

func NewWebsocketClient(conn *websocket.Conn, manager *Manager) *ClientHandeller {
	return &ClientHandeller{
		conn:    conn,
		manager: manager,
	}
}

func (client *ClientHandeller) SendToNseChan(stock *pb.Stock) {
	client.nseChan <- stock
}

func (client *ClientHandeller) ReadMessages() {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			client.manager.RemoveClient(client)
			break
		}

		var req clientRequest
		err = json.Unmarshal(msg, &req)
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
			client.manager.AddClient(client)
			logger.Printf("Client added to the list %v", len(client.manager.ClientsList))
		}
	}
}

func (m *Manager) SendStocks(stk *pb.Stock) {
	var i int = 1
		if len(m.ClientsList) > 0 {
			for client := range m.ClientsList {
				stkBytes, _ := json.Marshal(stk)
				client.conn.WriteMessage(websocket.TextMessage, stkBytes)
				i++
			}
		}


}
