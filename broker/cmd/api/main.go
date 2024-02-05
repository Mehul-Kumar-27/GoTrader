package main

import (
	"gotrader/broker/cmd/api/listner"
	logger "gotrader/logger"
	pb "gotrader/proto"

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
	listner.GetStocksForNse(client)
}
