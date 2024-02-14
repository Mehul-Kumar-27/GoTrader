package listner

import (
	"context"
	lg "gotrader/logger"
	pb "gotrader/proto"
	"gotrader/server/cmd/api/routes"

	"log"
)

var logger *log.Logger

func init() {
	logger = lg.CreateCustomLogger("listner")

}

func GetStocksForNse(client pb.StockServiceClient, manager *routes.Manager) {
	stream, err := client.GetStocks(context.Background(), &pb.ExchangeRequest{Exchange: "NSE"})

	if err != nil {
		logger.Fatalf("Error getting stocks: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err != nil {
			logger.Fatalf("Error receiving message: %v", err)
		}
		logger.Printf("Stock: %v price: %v", message.Name, message.Price)
		manager.SendStocks(message)
	}

}
