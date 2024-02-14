package main

import (
	"gotrader/server/cmd/api/listner"
	"gotrader/server/cmd/api/routes"
	"net/http"

	logger "gotrader/logger"
	pb "gotrader/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := logger.CreateCustomLogger("broker/api")
	conn, err := grpc.Dial("store:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Panicf("Failed to listen: %v", err)
	}
	defer conn.Close()

	client := pb.NewStockServiceClient(conn)
	logger.Println("Connected to store")
	logger.Println("Trying to start the server")
	router := mux.NewRouter()
	mangaer := routes.NewManager()
	go listner.GetStocksForNse(client, mangaer)
	router.HandleFunc("/api", mangaer.HandleRoutes).Methods("GET")
	logger.Println("There is websocket")
	router.HandleFunc("/ws", mangaer.HandleWebSocket)

	logger.Println("Server is running on port 8080")
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	er := server.ListenAndServe()
	if er != nil {
		logger.Panicf("Failed to listen: %v", er)
	}

}
