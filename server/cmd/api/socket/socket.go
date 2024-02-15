package socket

import (
	"encoding/json"
	lg "gotrader/logger"
	pb "gotrader/proto"
	"gotrader/server/cmd/api/listner"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var websockerUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("server/socket")
}

type sockerRequest struct {
	Exchange string `json:"exchange"`
}

type Manger struct{}

func NewManger() *Manger {
	return &Manger{}
}

func HandleWebSockerConnection(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Handling websocket connection")
	conn, err := websockerUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("Error upgrading to websocket: %v", err)
		return
	}
	stkChn := make(chan *pb.Stock)
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Printf("Error reading message: %v", err)
			break
		}
		logger.Printf("Received message: %v", string(msg))
		decoder := json.NewDecoder(strings.NewReader(string(msg)))
		var request sockerRequest
		err = decoder.Decode(&request)
		if err != nil {
			logger.Printf("Error decoding json: %v", err)
			break
		}
		if request.Exchange == "NSE" {
			go getData(w, stkChn)
			for s := range stkChn {
				conn.WriteMessage(websocket.TextMessage, []byte(s.String()))
			}
		}

	}

}

func (m *Manger) InitalRoute(w http.ResponseWriter, r *http.Request) {
	err := writeJson(w, http.StatusOK, jsonResponse{Error: false, Message: "Connected to socket server"})
	if err != nil {
		logger.Printf("Error writing json: %v", err)
		writeError(w, err)
	}
}

func getData(w http.ResponseWriter, stkChn chan *pb.Stock) {
	conn, err := grpc.Dial("store:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Printf("Failed to dial: %v", err)
		writeError(w, err)
	}

	defer conn.Close()
	client := pb.NewStockServiceClient(conn)

	listner.GetStocksForNse(client, stkChn)

}
