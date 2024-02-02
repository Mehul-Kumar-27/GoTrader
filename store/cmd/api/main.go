package main

import (
	lg "gotrader/logger"
	pb "gotrader/proto"
	"gotrader/store/cmd/api/stream"

	"gotrader/store/cmd/api/consumer"
	"log"
	"net"

	"google.golang.org/grpc"
)

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("store/api")
}

func main() {
	consumer.ConsumeMessages()
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Panicf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	st := stream.Stream{}
	pb.RegisterStockServiceServer(grpcServer, &st)
	logger.Println("Starting server on port :8080")

	if err := grpcServer.Serve(list); err != nil {
		logger.Panicf("Failed to serve: %v", err)
	}
}
